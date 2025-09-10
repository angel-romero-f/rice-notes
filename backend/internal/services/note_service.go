package services

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/angel-romero-f/rice-notes/internal/infra/storage"
	"github.com/angel-romero-f/rice-notes/internal/models"
	"github.com/angel-romero-f/rice-notes/internal/repository"
	"github.com/google/uuid"
)

const (
	// MaxFileSize is the maximum allowed file size (10MB)
	MaxFileSize = 10 * 1024 * 1024
	// AllowedContentType is the only allowed content type
	AllowedContentType = "application/pdf"
)

// NoteService handles note-related business logic
type NoteService struct {
	repo     repository.NoteRepository
	uploader storage.Uploader
}

// NewNoteService creates a new note service instance
func NewNoteService(repo repository.NoteRepository, uploader storage.Uploader) *NoteService {
	return &NoteService{
		repo:     repo,
		uploader: uploader,
	}
}

// CreateNote creates a new note by uploading a PDF file
func (s *NoteService) CreateNote(ctx context.Context, userEmail, title, courseID string, file multipart.File, header *multipart.FileHeader) (*models.NoteResponse, error) {
	slog.Info("Creating new note", "userEmail", userEmail, "title", title, "courseID", courseID, "fileName", header.Filename)

	// Validate inputs
	if err := s.validateCreateNoteRequest(userEmail, title, courseID, header); err != nil {
		slog.Warn("Invalid create note request", "error", err)
		return nil, err
	}

	// Generate UUID for the note
	noteID := uuid.New()

	// Create note model
	note := &models.Note{
		ID:          noteID,
		UserEmail:   userEmail,
		Title:       title,
		CourseID:    courseID,
		FileName:    header.Filename,
		FileSize:    header.Size,
		ContentType: AllowedContentType,
		FilePath:    storage.GenerateFileKey(userEmail, noteID.String(), header.Filename),
	}

	// Upload file to S3
	if err := s.uploader.Upload(ctx, note.FilePath, file, note.ContentType, note.FileSize); err != nil {
		slog.Error("Failed to upload file to S3", "error", err, "noteID", noteID)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Save note to database
	if err := s.repo.CreateNote(ctx, note); err != nil {
		// Try to clean up uploaded file on database error
		if deleteErr := s.uploader.Delete(ctx, note.FilePath); deleteErr != nil {
			slog.Error("Failed to cleanup file after database error", "deleteError", deleteErr, "noteID", noteID)
		}
		slog.Error("Failed to save note to database", "error", err, "noteID", noteID)
		return nil, fmt.Errorf("failed to save note: %w", err)
	}

	// Return response
	response := &models.NoteResponse{
		ID:          note.ID,
		Title:       note.Title,
		CourseID:    note.CourseID,
		FileName:    note.FileName,
		FileSize:    note.FileSize,
		ContentType: note.ContentType,
		UploadedAt:  note.UploadedAt,
	}

	slog.Info("Note created successfully", "noteID", noteID, "userEmail", userEmail)
	return response, nil
}

// GetNoteByID retrieves a note by its ID
func (s *NoteService) GetNoteByID(ctx context.Context, noteID uuid.UUID, userEmail string) (*models.Note, error) {
	note, err := s.repo.GetNoteByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	// Ensure the note belongs to the requesting user
	if note.UserEmail != userEmail {
		slog.Warn("User attempted to access note they don't own", 
			"userEmail", userEmail, "noteOwner", note.UserEmail, "noteID", noteID)
		return nil, fmt.Errorf("note not found")
	}

	return note, nil
}

// GetUserNotes retrieves notes for a user with optional course filtering
func (s *NoteService) GetUserNotes(ctx context.Context, userEmail, courseID string, limit, offset int) ([]*models.Note, error) {
	// Apply reasonable limits
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var notes []*models.Note
	var err error

	if courseID != "" {
		notes, err = s.repo.GetNotesByCourse(ctx, userEmail, courseID, limit, offset)
	} else {
		notes, err = s.repo.GetNotesByUser(ctx, userEmail, limit, offset)
	}

	if err != nil {
		slog.Error("Failed to get user notes", "error", err, "userEmail", userEmail)
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}

	return notes, nil
}

// DeleteNote deletes a note and its associated file
func (s *NoteService) DeleteNote(ctx context.Context, noteID uuid.UUID, userEmail string) error {
	// First get the note to check ownership and get file path
	note, err := s.GetNoteByID(ctx, noteID, userEmail)
	if err != nil {
		return err
	}

	// Delete from database first
	if err := s.repo.DeleteNote(ctx, noteID, userEmail); err != nil {
		slog.Error("Failed to delete note from database", "error", err, "noteID", noteID)
		return fmt.Errorf("failed to delete note: %w", err)
	}

	// Delete file from S3 (best effort - don't fail if this fails)
	if err := s.uploader.Delete(ctx, note.FilePath); err != nil {
		slog.Error("Failed to delete file from S3", "error", err, "noteID", noteID, "filePath", note.FilePath)
		// Don't return error here - the database deletion was successful
	}

	slog.Info("Note deleted successfully", "noteID", noteID, "userEmail", userEmail)
	return nil
}

// validateCreateNoteRequest validates the request parameters
func (s *NoteService) validateCreateNoteRequest(userEmail, title, courseID string, header *multipart.FileHeader) error {
	if userEmail == "" {
		return fmt.Errorf("user email is required")
	}

	if title == "" {
		return fmt.Errorf("title is required")
	}

	if len(title) > 255 {
		return fmt.Errorf("title must be 255 characters or less")
	}

	if courseID == "" {
		return fmt.Errorf("course ID is required")
	}

	if len(courseID) > 50 {
		return fmt.Errorf("course ID must be 50 characters or less")
	}

	if header == nil {
		return fmt.Errorf("file is required")
	}

	if header.Size == 0 {
		return fmt.Errorf("file cannot be empty")
	}

	if header.Size > MaxFileSize {
		return fmt.Errorf("file size must be less than %d bytes", MaxFileSize)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		return fmt.Errorf("only PDF files are allowed")
	}

	return nil
}

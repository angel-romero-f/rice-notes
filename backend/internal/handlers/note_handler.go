package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/angel-romero-f/rice-notes/internal/middleware"
	"github.com/angel-romero-f/rice-notes/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// NoteService defines the business logic interface for note operations
type NoteService interface {
	CreateNote(ctx context.Context, userEmail, title, courseID string, file multipart.File, header *multipart.FileHeader) (*models.NoteResponse, error)
	GetNoteByID(ctx context.Context, noteID uuid.UUID, userEmail string) (*models.Note, error)
	GetUserNotes(ctx context.Context, userEmail, courseID string, limit, offset int) ([]*models.Note, error)
	DeleteNote(ctx context.Context, noteID uuid.UUID, userEmail string) error
}

// NoteHandler handles HTTP requests for note operations
type NoteHandler struct {
	service NoteService
}

// NewNoteHandler creates a new note handler instance
func NewNoteHandler(service NoteService) *NoteHandler {
	return &NoteHandler{
		service: service,
	}
}

// CreateNote handles POST /api/notes - uploads a PDF file and creates a note
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	// Get user from JWT context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("User not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (32MB max memory)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		slog.Error("Failed to parse multipart form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract form fields
	title := r.FormValue("title")
	courseID := r.FormValue("course_id")

	if title == "" || courseID == "" {
		http.Error(w, "Title and course_id are required", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to get uploaded file", "error", err)
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create note
	response, err := h.service.CreateNote(r.Context(), user.Email, title, courseID, file, header)
	if err != nil {
		slog.Error("Failed to create note", "error", err, "userEmail", user.Email)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Note created successfully", "noteID", response.ID, "userEmail", user.Email)
}

// GetNotes handles GET /api/notes - retrieves notes for the authenticated user
func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	// Get user from JWT context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("User not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	courseID := r.URL.Query().Get("course_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default pagination values
	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get notes
	notes, err := h.service.GetUserNotes(r.Context(), user.Email, courseID, limit, offset)
	if err != nil {
		slog.Error("Failed to get notes", "error", err, "userEmail", user.Email)
		http.Error(w, "Failed to get notes", http.StatusInternalServerError)
		return
	}

	// Return notes
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug("Notes retrieved", "userEmail", user.Email, "count", len(notes))
}

// GetNote handles GET /api/notes/{id} - retrieves a specific note (returns 302 redirect to S3)
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	// Get user from JWT context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("User not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse note ID from URL
	noteIDStr := chi.URLParam(r, "id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		slog.Error("Invalid note ID", "noteID", noteIDStr, "error", err)
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Get note
	note, err := h.service.GetNoteByID(r.Context(), noteID, user.Email)
	if err != nil {
		slog.Error("Failed to get note", "error", err, "noteID", noteID, "userEmail", user.Email)
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// For now, return note metadata. In the future, this could redirect to a presigned S3 URL
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(note); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Debug("Note retrieved", "noteID", noteID, "userEmail", user.Email)
}

// DeleteNote handles DELETE /api/notes/{id} - deletes a note
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	// Get user from JWT context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		slog.Error("User not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse note ID from URL
	noteIDStr := chi.URLParam(r, "id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		slog.Error("Invalid note ID", "noteID", noteIDStr, "error", err)
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Delete note
	if err := h.service.DeleteNote(r.Context(), noteID, user.Email); err != nil {
		slog.Error("Failed to delete note", "error", err, "noteID", noteID, "userEmail", user.Email)
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
	slog.Info("Note deleted", "noteID", noteID, "userEmail", user.Email)
}

// Welcome handles GET / - returns welcome message (keeping for backward compatibility)
func (h *NoteHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Welcome to Rice Notes!",
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode welcome response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

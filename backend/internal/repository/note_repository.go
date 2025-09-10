// Package repository provides database access layer for notes
package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/angel-romero-f/rice-notes/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NoteRepository defines the interface for note database operations
type NoteRepository interface {
	CreateNote(ctx context.Context, note *models.Note) error
	GetNoteByID(ctx context.Context, id uuid.UUID) (*models.Note, error)
	GetNotesByUser(ctx context.Context, userEmail string, limit, offset int) ([]*models.Note, error)
	GetNotesByCourse(ctx context.Context, userEmail, courseID string, limit, offset int) ([]*models.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID, userEmail string) error
}

// PostgresNoteRepository implements NoteRepository using PostgreSQL
type PostgresNoteRepository struct {
	db *pgxpool.Pool
}

// NewPostgresNoteRepository creates a new PostgreSQL-based note repository
func NewPostgresNoteRepository(db *pgxpool.Pool) *PostgresNoteRepository {
	return &PostgresNoteRepository{
		db: db,
	}
}

// CreateNote creates a new note in the database
func (r *PostgresNoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO notes (id, user_email, title, course_id, file_name, file_path, file_size, content_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING uploaded_at, updated_at`
	
	err := r.db.QueryRow(ctx, query,
		note.ID,
		note.UserEmail,
		note.Title,
		note.CourseID,
		note.FileName,
		note.FilePath,
		note.FileSize,
		note.ContentType,
	).Scan(&note.UploadedAt, &note.UpdatedAt)

	if err != nil {
		slog.Error("Failed to create note", "error", err, "noteID", note.ID)
		return fmt.Errorf("failed to create note: %w", err)
	}

	slog.Info("Note created successfully", "noteID", note.ID, "userEmail", note.UserEmail)
	return nil
}

// GetNoteByID retrieves a note by its ID
func (r *PostgresNoteRepository) GetNoteByID(ctx context.Context, id uuid.UUID) (*models.Note, error) {
	query := `
		SELECT id, user_email, title, course_id, file_name, file_path, file_size, 
			   content_type, uploaded_at, updated_at
		FROM notes 
		WHERE id = $1`
	
	note := &models.Note{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&note.ID,
		&note.UserEmail,
		&note.Title,
		&note.CourseID,
		&note.FileName,
		&note.FilePath,
		&note.FileSize,
		&note.ContentType,
		&note.UploadedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Debug("Note not found", "noteID", id)
			return nil, fmt.Errorf("note not found: %s", id)
		}
		slog.Error("Failed to get note by ID", "error", err, "noteID", id)
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	slog.Debug("Note retrieved successfully", "noteID", id)
	return note, nil
}

// GetNotesByUser retrieves notes for a specific user with pagination
func (r *PostgresNoteRepository) GetNotesByUser(ctx context.Context, userEmail string, limit, offset int) ([]*models.Note, error) {
	query := `
		SELECT id, user_email, title, course_id, file_name, file_path, file_size, 
			   content_type, uploaded_at, updated_at
		FROM notes 
		WHERE user_email = $1
		ORDER BY uploaded_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(ctx, query, userEmail, limit, offset)
	if err != nil {
		slog.Error("Failed to query notes by user", "error", err, "userEmail", userEmail)
		return nil, fmt.Errorf("failed to get notes for user: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(
			&note.ID,
			&note.UserEmail,
			&note.Title,
			&note.CourseID,
			&note.FileName,
			&note.FilePath,
			&note.FileSize,
			&note.ContentType,
			&note.UploadedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan note", "error", err)
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating rows", "error", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	slog.Debug("Notes retrieved for user", "userEmail", userEmail, "count", len(notes))
	return notes, nil
}

// GetNotesByCourse retrieves notes for a specific user and course with pagination
func (r *PostgresNoteRepository) GetNotesByCourse(ctx context.Context, userEmail, courseID string, limit, offset int) ([]*models.Note, error) {
	query := `
		SELECT id, user_email, title, course_id, file_name, file_path, file_size, 
			   content_type, uploaded_at, updated_at
		FROM notes 
		WHERE user_email = $1 AND course_id = $2
		ORDER BY uploaded_at DESC
		LIMIT $3 OFFSET $4`
	
	rows, err := r.db.Query(ctx, query, userEmail, courseID, limit, offset)
	if err != nil {
		slog.Error("Failed to query notes by course", "error", err, 
			"userEmail", userEmail, "courseID", courseID)
		return nil, fmt.Errorf("failed to get notes for course: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(
			&note.ID,
			&note.UserEmail,
			&note.Title,
			&note.CourseID,
			&note.FileName,
			&note.FilePath,
			&note.FileSize,
			&note.ContentType,
			&note.UploadedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan note", "error", err)
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating rows", "error", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	slog.Debug("Notes retrieved for course", "userEmail", userEmail, "courseID", courseID, "count", len(notes))
	return notes, nil
}

// DeleteNote deletes a note (only if it belongs to the specified user)
func (r *PostgresNoteRepository) DeleteNote(ctx context.Context, id uuid.UUID, userEmail string) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_email = $2`
	
	result, err := r.db.Exec(ctx, query, id, userEmail)
	if err != nil {
		slog.Error("Failed to delete note", "error", err, "noteID", id, "userEmail", userEmail)
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		slog.Warn("Note not found or not owned by user", "noteID", id, "userEmail", userEmail)
		return fmt.Errorf("note not found or not owned by user")
	}

	slog.Info("Note deleted successfully", "noteID", id, "userEmail", userEmail)
	return nil
}
package models

import (
	"time"

	"github.com/google/uuid"
)

// Note represents a PDF note uploaded by a user
type Note struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserEmail   string    `json:"user_email" db:"user_email"`
	Title       string    `json:"title" db:"title"`
	CourseID    string    `json:"course_id" db:"course_id"`
	FileName    string    `json:"file_name" db:"file_name"`
	FilePath    string    `json:"file_path" db:"file_path"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	ContentType string    `json:"content_type" db:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateNoteRequest represents the request payload for creating a new note
type CreateNoteRequest struct {
	Title    string `json:"title" form:"title"`
	CourseID string `json:"course_id" form:"course_id"`
	// File will be handled separately in multipart form
}

// NoteResponse represents the response when returning note information
type NoteResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	CourseID    string    `json:"course_id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

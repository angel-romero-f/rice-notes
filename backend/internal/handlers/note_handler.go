package handlers

import (
	"fmt"
	"net/http"
)
// NoteService defines the business logic for note operations 
type NoteService interface {
	UploadNote() string
}
type NoteHandler struct {
	s NoteService
}

// NewNoteHandler returns a concrete implementation of NoteService interface. Has handlers
// for CreateNote.
func NewNoteHandler(s NoteService) *NoteHandler {
	return &NoteHandler{s: s}
}

func (n *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	response := n.s.UploadNote()
	w.Write([]byte(fmt.Sprintf("%s", response)))
}

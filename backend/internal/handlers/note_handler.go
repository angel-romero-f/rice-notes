package handlers

import (
	"fmt"
	"net/http"
)

type NoteService interface {
	UploadNote() string
}
type NoteHandler struct {
	s NoteService
}

func NewNoteHandler(s NoteService) *NoteHandler {
	return &NoteHandler{s: s}
}

func (n *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	response := n.s.UploadNote()
	w.Write([]byte(fmt.Sprintf("%s", response)))
}

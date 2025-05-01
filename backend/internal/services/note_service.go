package services

import "fmt"

type NoteUploader interface {
	UploadNote(text string) error
}

type noteUploader struct {
}

func NewNoteService() NoteUploader {
	return &noteUploader{}
}

func (n *noteUploader) UploadNote(text string) error {
	fmt.Println("recieved note! %s", text)
	return nil
}

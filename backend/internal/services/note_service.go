package services

type noteService struct {
}

// Factory method that returns a new NoteService instnace.
func NewNoteService() *noteService {
	return &noteService{}
}

// UploadNote reutrns the string to be logged to the user after making a request.
func (s *noteService) UploadNote() string {
	return "Welcome to rice notes!"
}

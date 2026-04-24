package models

//easyjson:json
type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required"`
}

//easyjson:json
type NoteList []Note

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error"`
}
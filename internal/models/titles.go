package models

type (
	Title struct {
		ID         int      `json:"id" DB:"id"`
		Name       string   `json:"name" DB:"name"`
		CategoryID string   `json:"category_id" DB:"category_id"`
		Chapters   Chapters `json:"chapters" DB:"chapters"`
	}
	Chapters struct {
		Number int
		pages  Pages
	}
	Pages struct {
		Img    string
		Number int
	}
)

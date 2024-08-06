package models

type (
	Publisher struct {
		ID       int    `json:"id" DB:"id"`
		Name     string `json:"name" DB:"name"`
		Login_ID Login_ID
		Library  library_of_released_titles
	}
	Publisherfilter struct {
		Query *string `json:"query"`
	}
	library_of_released_titles struct {
		Publisher_ID int
		titles       []Title
	}
)

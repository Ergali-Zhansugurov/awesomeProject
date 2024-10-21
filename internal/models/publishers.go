package models

type (
	Publisher struct {
		ID      int                        `json:"id" db:"id"`
		Name    string                     `json:"name" db:"name"`
		LoginID LoginID                    `json:"login_id"`
		Library library_of_released_titles `json:"library"`
	}
	library_of_released_titles struct {
		Publisher_ID int
		titles       []Title
	}
)

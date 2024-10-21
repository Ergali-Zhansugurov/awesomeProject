package models

type (
	Category struct {
		Name string `json:"name" DB:"name"`
		id   int
	}
)

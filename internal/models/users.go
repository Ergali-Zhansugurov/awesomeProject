package models

type (
	User struct {
		ID       int    `json:"id" DB:"id"`
		Name     string `json:"name" DB:"name"`
		Login_ID Login_ID
		Library  User_library
	}
	User_library struct {
		read_now     []Read_now
		already_read []Title
		user_ID      int
	}
	Read_now struct {
		Title
		Chapters int `json:"chapters" DB:"chapters"`
	}
	Login_ID struct {
		Username string
		password string
	}
	UserFilter struct {
		Query *string `json:"query"`
	}
)

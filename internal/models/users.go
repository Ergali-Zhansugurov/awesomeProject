package models

type User struct {
	ID      int         `json:"id" db:"id"`
	Name    string      `json:"name" db:"name"`
	LoginID LoginID     `json:"login_id"`
	Library UserLibrary `json:"library"`
}

type UserLibrary struct {
	ReadNow     []ReadNow `json:"read_now"`
	AlreadyRead []Title   `json:"already_read"`
	UserID      int       `json:"user_id" db:"user_id"`
}

type ReadNow struct {
	Title    `json:"title"`
	Chapters int `json:"chapters" db:"chapters"`
}

type LoginID struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

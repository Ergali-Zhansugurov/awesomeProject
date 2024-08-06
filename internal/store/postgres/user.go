package postgres

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/store"
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	data map[int]*models.Title
	conn *sqlx.DB
	mu   *sync.RWMutex
}

func (db *DB) User() store.UserRepository {
	if db.User() != nil {
		panic(db.User())
	}
	return NewUsersRepository(db.conn)
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UserRepository {
	return &UsersRepository{conn: conn}
}

func (u *UsersRepository) Get(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	basicQuery := "SELECT *FROM Users"
	searchQuery := ""
	if filter.Query != nil {
		basicQuery += " WHERE name ilike '%$1%'" + *filter.Query + "%''"
		searchQuery = *filter.Query
	}
	User := make([]*models.User, 0)
	if err := u.conn.Select(&User, basicQuery, searchQuery); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return User, nil
}
func (u *UsersRepository) ByID(ctx context.Context, id int) (*models.User, error) {
	User := new(models.User)
	if err := u.conn.Get(User, "SELECT id , name FROM Users WHERE id=$1", id); err != nil {
		fmt.Printf(":%v :%s", nil, err)
	}
	return User, nil
}
func (u *UsersRepository) Create(ctx context.Context, User *models.User) error {
	_, err := u.conn.Exec("INSERT INTO Users(name) VALUES ($1)", User.Name)
	if err != nil {
		return fmt.Errorf("unknow err: %s", err)
	}
	return nil
}
func (u *UsersRepository) Delete(ctx context.Context, id int) error {
	_, err := u.conn.Exec("DELETE FROM users WHERE User_id=$1", id)
	if err != nil {
		panic(err)
	}
	return nil
}
func (u *UsersRepository) Update(ctx context.Context, User *models.User) error {
	_, err := u.conn.Exec("INSERT INTO Users(name) VALUES ($1)", User.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}

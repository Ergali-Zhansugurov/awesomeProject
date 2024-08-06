package postgres

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/store"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Publisher() store.PublisherRepository {
	if db.Publisher() != nil {
		panic(db.Publisher())
	}
	return NewPublisherRepository(db.conn)
}

type PublisherRepository struct {
	conn *sqlx.DB
}

func NewPublisherRepository(conn *sqlx.DB) store.PublisherRepository {
	return &PublisherRepository{conn: conn}
}

func (u *PublisherRepository) Get(ctx context.Context, filter *models.Publisherfilter) ([]*models.Publisher, error) {
	basicQuery := "SELECT *FROM Users"
	searchQuery := ""
	if filter.Query != nil {
		basicQuery += " WHERE name ilike '%$1%'" + *filter.Query + "%''"
		searchQuery = *filter.Query
	}
	Publisher := make([]*models.Publisher, 0)
	if err := u.conn.Select(&Publisher, basicQuery, searchQuery); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return Publisher, nil
}
func (u *PublisherRepository) ByID(ctx context.Context, id int) (*models.Publisher, error) {
	Publisher := new(models.Publisher)
	if err := u.conn.Get(Publisher, "SELECT id , name FROM Publishers WHERE id=$1", id); err != nil {
		fmt.Printf(":%v :%s", nil, err)
	}
	return Publisher, nil
}
func (u *PublisherRepository) Create(ctx context.Context, Publisher *models.Publisher) error {
	_, err := u.conn.Exec("INSERT INTO Publishers(name) VALUES ($1)", Publisher.Name)
	if err != nil {
		return fmt.Errorf("unknow err: %s", err)
	}
	return nil
}
func (u *PublisherRepository) Delete(ctx context.Context, id int) error {
	_, err := u.conn.Exec("DELETE FROM Publishers WHERE User_id=$1", id)
	if err != nil {
		panic(err)
	}
	return nil
}
func (u *PublisherRepository) Update(ctx context.Context, Publisher *models.Publisher) error {
	_, err := u.conn.Exec("INSERT INTO Publishers(name) VALUES ($1)", Publisher.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}

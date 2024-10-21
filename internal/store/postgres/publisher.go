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
func (p *PublisherRepository) AddTitleToLibrary(ctx context.Context, userID int, title *models.Title) error {
	_, err := p.conn.Exec("INSERT INTO user_library (user_id, title_id) VALUES ($1, $2)", userID, title.ID)
	return err
}
func (p *PublisherRepository) RemoveTitleFromLibrary(ctx context.Context, userID int, titleID int) error {
	_, err := p.conn.Exec("DELETE FROM user_library WHERE user_id = $1 AND title_id = $2", userID, titleID)
	return err
}

func (p *PublisherRepository) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User
	err := p.conn.Get(&user, "SELECT id, name FROM users WHERE username = $1 AND password = $2", username, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	//todo
	return &user, nil
}
func (p *PublisherRepository) Get(ctx context.Context, filter *models.Filter) ([]*models.Publisher, error) {
	basicQuery := "SELECT *FROM Users"
	searchQuery := ""
	if filter.Query != nil {
		basicQuery += " WHERE name ilike '%$1%'" + *filter.Query + "%''"
		searchQuery = *filter.Query
	}
	Publisher := make([]*models.Publisher, 0)
	if err := p.conn.Select(&Publisher, basicQuery, searchQuery); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return Publisher, nil
}
func (p *PublisherRepository) ByID(ctx context.Context, id int) (*models.Publisher, error) {
	Publisher := new(models.Publisher)
	if err := p.conn.Get(Publisher, "SELECT id , name FROM Publishers WHERE id=$1", id); err != nil {
		fmt.Printf(":%v :%s", nil, err)
	}
	return Publisher, nil
}
func (p *PublisherRepository) Create(ctx context.Context, Publisher *models.Publisher) error {
	_, err := p.conn.Exec("INSERT INTO Publishers(name) VALUES ($1)", Publisher.Name)
	if err != nil {
		return fmt.Errorf("unknow err: %s", err)
	}
	return nil
}
func (p *PublisherRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec("DELETE FROM Publishers WHERE User_id=$1", id)
	if err != nil {
		panic(err)
	}
	return nil
}
func (p *PublisherRepository) Update(ctx context.Context, Publisher *models.Publisher) error {
	_, err := p.conn.Exec("INSERT INTO Publishers(name) VALUES ($1)", Publisher.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}

package postgres

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/store"
	"context"
	"fmt"

	"sync"

	"github.com/jmoiron/sqlx"
)

type TitlesRepo struct {
	data map[int]*models.Title
	conn *sqlx.DB
	mu   *sync.RWMutex
}

func (db *DB) Title() store.TitleRepository {
	if db.Title() != nil {
		panic(db.Title())
	}
	return NewTitlesRepository(db.conn)
}

type TitlesRepository struct {
	conn *sqlx.DB
}

func (c *TitlesRepository) ByCategoryId(ctx context.Context, category_id int) (*models.Title, error) {
	Title := new(models.Title)
	if err := c.conn.Get(Title, "SELECT id , name FROM Title WHERE category_id=$1", category_id); err != nil {
		fmt.Printf(":%v :%s", nil, err)
	}
	return Title, nil
}

func NewTitlesRepository(conn *sqlx.DB) store.TitleRepository {
	return &TitlesRepository{conn: conn}
}
func (c TitlesRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM Titles WHERE user_id=$1", id)
	if err != nil {
		panic(err)
	}
	return nil
}
func (c TitlesRepository) Create(ctx context.Context, Title *models.Title) error {
	_, err := c.conn.Exec("INSERT INTO Titles(name) VALUES ($1)", Title.Name)
	if err != nil {
		return fmt.Errorf("unknow err: %s", err)
	}
	return nil
}
func (c TitlesRepository) All(ctx context.Context, filter *models.Titlesfilter) ([]*models.Title, error) {
	basicQuery := "SELECT *FROM Titles"
	searchQuery := ""
	if filter.Query != nil {
		basicQuery += " WHERE name ilike '%$1%'" + *filter.Query + "%''"
		searchQuery = *filter.Query
	}
	Title := make([]*models.Title, 0)
	if err := c.conn.Select(&Title, basicQuery, searchQuery); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return Title, nil
}
func (c TitlesRepository) ByID(ctx context.Context, id int) (*models.Title, error) {
	Title := new(models.Title)
	if err := c.conn.Get(Title, "SELECT id , name FROM Title WHERE id=$1", id); err != nil {
		fmt.Printf(":%v :%s", nil, err)
	}
	return Title, nil
}
func (c TitlesRepository) Update(ctx context.Context, Title *models.Title) error {
	_, err := c.conn.Exec("INSERT INTO Titles(name) VALUES ($1)", Title.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}

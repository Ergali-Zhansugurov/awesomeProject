package postgres

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/store"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Category() store.CategoryRepository {
	if db.Category() != nil {
		panic(db.Category())
	}
	return NewCategoryRepository(db.conn)
}

type CategoriesRepository struct {
	conn *sqlx.DB
}

func NewCategoryRepository(conn *sqlx.DB) store.CategoryRepository {
	return &CategoriesRepository{conn: conn}
}

func (c CategoriesRepository) Create(ctx context.Context, category *models.Category) error {
	_, err := c.conn.Exec("INSERT INTO categories(name) VALUES ($1)", category.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}
func (c CategoriesRepository) Update(ctx context.Context, category *models.Category) error {
	_, err := c.conn.Exec("INSERT INTO categories(name) VALUES ($1)", category.Name)
	if err != nil {
		return fmt.Errorf("Unknow err:%S", err)
	}
	return nil
}
func (c CategoriesRepository) Get(ctx context.Context, filter *models.Categoryesfilter) ([]*models.Category, error) {
	basicQuery := "SELECT *FROM categories"
	searchQuery := ""
	if filter.Query != nil {
		basicQuery += " WHERE name ilike '%$1%'" + *filter.Query + "%''"
		searchQuery = *filter.Query
	}
	categories := make([]*models.Category, 0)
	if err := c.conn.Select(&categories, basicQuery, searchQuery); err != nil {
		return nil, fmt.Errorf("%S", err)
	}
	return categories, nil
}
func (c CategoriesRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM Categories WHERE user_id=$1", id)
	if err != nil {
		panic(err)
	}
	return nil
}

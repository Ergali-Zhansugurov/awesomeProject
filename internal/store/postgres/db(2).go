package postgres

import (
	"awesomeProject/internal/store"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn       *sqlx.DB
	Titles     store.TitleRepository
	Categorys  store.CategoryRepository
	Publishers store.PublisherRepository
	Users      store.UserRepository
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Connect(urlExample string) error {
	conn, err := sqlx.Connect("pgx", urlExample)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}
	db.conn = conn
	return nil
}
func NewDB() store.Store {
	return &DB{}
}

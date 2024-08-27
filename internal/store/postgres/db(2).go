package postgres

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/store"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const PG_URL = "postgres://postgres:%v@%v:%v/%v"

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

func (db *DB) Connect(cfg *config.StorageConfig) error {

	url := fmt.Sprintf(PG_URL, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	conn, err := sqlx.Connect("pgx", url)
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

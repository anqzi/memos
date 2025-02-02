package store

import (
	"database/sql"
	"log"

	"github.com/kehuay/aimemos/store/db"
)

type Store struct {
	db *sql.DB
}

func NewStore(config db.DBConfig) (Store, error) {
	_db, err := db.NewPostgresDB(config)
	if err != nil {
		return Store{}, err
	}
	log.Println("数据库连接成功")
	return Store{
		db: _db,
	}, nil
}

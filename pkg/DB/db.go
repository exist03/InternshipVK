package DB

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDB(dst string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dst)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}

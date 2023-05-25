package repository

import (
	"database/sql"
	"log"
)

type Repository struct {
	DB *sql.DB
}

func New(DB *sql.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Insert(username, service, login, password interface{}) error {
	stmt := `INSERT INTO Services (username, service, login, password) VALUES(?, ?, ?, ?)`

	_, err := r.DB.Exec(stmt, username, service, login, password)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) Delete(username, service string) error {
	stmt := `DELETE FROM Services WHERE username=? AND service=?`

	_, err := r.DB.Exec(stmt, username, service)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) Get(username, service string) (string, string) {
	var login, password string
	stmt := `SELECT login, password FROM Services WHERE username=? AND service=?`
	row := r.DB.QueryRow(stmt, username, service)
	err := row.Scan(&login, &password)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return login, password
}
func (r *Repository) GetList(userID string) ([]string, error) {
	var services []string
	stmt := `SELECT service FROM Services WHERE username=?`
	rows, err := r.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var service string
	for rows.Next() {
		err := rows.Scan(&service)
		if err != nil {
			log.Print(err)
		}
		services = append(services, service)
	}
	services = append(services, "/cancel")
	return services, nil
}

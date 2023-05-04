package mysql

import (
	"database/sql"
	"log"
)

type FormModel struct {
	DB *sql.DB
}

func (m *FormModel) Insert(username, service, login, password interface{}) error {
	stmt := `INSERT INTO Services (username, service, login, password) VALUES(?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt, username, service, login, password)
	if err != nil {
		return err
	}
	return nil
}

func (m *FormModel) Delete(username, service string) error {
	stmt := `DELETE FROM Services WHERE username=? AND service=?`

	_, err := m.DB.Exec(stmt, username, service)
	if err != nil {
		return err
	}
	return nil
}
func (m *FormModel) Get(username, service string) (string, string) {
	var login, password string
	stmt := `SELECT login, password FROM Services WHERE username=? AND service=?`
	row := m.DB.QueryRow(stmt, username, service)
	err := row.Scan(&login, &password)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return login, password
}

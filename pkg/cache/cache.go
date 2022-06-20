package cache

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	os.Remove("./cache.db")

	db, err := sql.Open("sqlite3", "./cache.db")
	if err != nil {
		return nil, err
	}

	sqlCreate := `create table cache(id varchar(30) primary key, data blob)`

	_, err = db.Exec(sqlCreate)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertCache(db *sql.DB, id string, data []byte) error {
	stmt, err := db.Prepare("insert into cache values (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, data)
	if err != nil {
		return err
	}

	return nil
}

func GetCache(db *sql.DB, id string) ([]byte, error) {
	stmt, err := db.Prepare("select data from cache where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var data []byte
	err = stmt.QueryRow(id).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

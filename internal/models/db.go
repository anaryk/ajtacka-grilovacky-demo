package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func InitDB(dataSourceName string, dbName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS` + dbName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`USE ` + dbName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS alkoholik (
            id INT AUTO_INCREMENT PRIMARY KEY,
            jmeno VARCHAR(255) NOT NULL,
            fotka MEDIUMTEXT,
            piva INT DEFAULT 0,
            tvrdy INT DEFAULT 0,
            nealko INT DEFAULT 0,
            last_drink DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) GetAlkoholikByID(id int) (*Alkoholik, error) {
	var alkoholik Alkoholik
	err := db.QueryRow("SELECT id, jmeno, fotka FROM alkoholik WHERE id = ?", id).Scan(&alkoholik.ID, &alkoholik.Jmeno, &alkoholik.Fotka)
	if err != nil {
		return nil, err
	}
	return &alkoholik, nil
}

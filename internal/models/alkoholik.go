package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Alkoholik struct {
	ID        int
	Jmeno     string
	Fotka     string
	Piva      int
	Tvrdy     int
	Nealko    int
	LastDrink time.Time
}

func (db *DB) AddDrink(alkoholikID int, drinkType string) error {
	var column string
	switch drinkType {
	case "pivo":
		column = "piva"
	case "tvrdy":
		column = "tvrdy"
	case "nealko":
		column = "nealko"
	default:
		return fmt.Errorf("unknown drink type: %s", drinkType)
	}

	var drinkCount int
	err := db.QueryRow("SELECT "+column+" FROM alkoholik WHERE id = ?", alkoholikID).Scan(&drinkCount)
	if err != nil {
		return err
	}

	var query string
	if drinkCount > 0 {
		query = fmt.Sprintf("UPDATE alkoholik SET %s = %s + 1, last_drink = CURRENT_TIMESTAMP WHERE id = ?", column, column)
	} else {
		query = fmt.Sprintf("UPDATE alkoholik SET %s = %s + 1 WHERE id = ?", column, column)
	}

	_, err = db.Exec(query, alkoholikID)
	return err
}

func (db *DB) GetLastDrinkTime(alkoholikID int) (time.Time, error) {
	var lastDrink sql.NullString
	err := db.QueryRow("SELECT last_drink FROM alkoholik WHERE id = ?", alkoholikID).Scan(&lastDrink)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	if !lastDrink.Valid {
		return time.Time{}, nil
	}

	lastDrinkTime, err := time.Parse("2006-01-02 15:04:05", lastDrink.String)
	if err != nil {
		return time.Time{}, err
	}
	return lastDrinkTime, nil
}

func (db *DB) GetAllAlkoholici() ([]Alkoholik, error) {
	var alkoholici []Alkoholik
	rows, err := db.Query("SELECT id, jmeno, fotka, piva, tvrdy, nealko, last_drink FROM alkoholik")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Alkoholik
		var lastDrinkStr []byte
		err := rows.Scan(&a.ID, &a.Jmeno, &a.Fotka, &a.Piva, &a.Tvrdy, &a.Nealko, &lastDrinkStr)
		if err != nil {
			return nil, err
		}
		a.LastDrink, _ = time.Parse("2006-01-02 15:04:05", string(lastDrinkStr))
		alkoholici = append(alkoholici, a)
	}
	return alkoholici, nil
}

func (db *DB) CreateAlkoholik(jmeno, fotka string) (int, error) {
	result, err := db.Exec("INSERT INTO alkoholik (jmeno, fotka) VALUES (?, ?)", jmeno, fotka)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

package database

import (
	"database/sql"

	"github.com/ahui2016/ipelago/stmt"
)

type TX interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// getText1 gets one text value from the database.
func getText1(tx TX, query string, args ...interface{}) (text string, err error) {
	row := tx.QueryRow(query, args...)
	err = row.Scan(&text)
	return
}

// getInt1 gets one number value from the database.
func getInt1(tx TX, query string, arg ...interface{}) (n int64, err error) {
	row := tx.QueryRow(query, arg...)
	err = row.Scan(&n)
	return
}

type Row interface {
	Scan(...interface{}) error
}

func scanIsland(row Row) (island Island, err error) {
	err = row.Scan(
		&island.ID,
		&island.Name,
		&island.Avatar,
		&island.Email,
		&island.Link,
		&island.Address,
		&island.Note,
	)
	return
}

func getIslandByID(tx TX, id string) (island Island, err error) {
	row := tx.QueryRow(stmt.GetIslandByID, id)
	if island, err = scanIsland(row); err != nil {
		return
	}
	island.Message, err = getText1(tx, stmt.GetLastMessage, id)
	return
}

func updateIsland(tx TX, island Island) error {
	_, err := tx.Exec(
		stmt.UpdateIsland,
		island.Name,
		island.Avatar,
		island.Email,
		island.Link,
		island.Address,
		island.Note,
		island.ID,
	)
	return err
}

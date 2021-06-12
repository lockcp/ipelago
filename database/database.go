package database

import (
	"database/sql"

	"github.com/ahui2016/ipelago/model"
	"github.com/ahui2016/ipelago/stmt"
	_ "github.com/mattn/go-sqlite3"
)

const MyID = "My-Island-ID"

type (
	Island = model.Island
)

type DB struct {
	Path string
	DB   *sql.DB
}

func (db *DB) Exec(query string, args ...interface{}) (err error) {
	_, err = db.DB.Exec(query, args...)
	return
}

func (db *DB) Open(dbPath string) (err error) {
	if db.DB, err = sql.Open("sqlite3", dbPath+"?_fk=1"); err != nil {
		return
	}
	db.Path = dbPath
	return db.Exec(stmt.CreateTables)
}

func (db *DB) MyIsland() (Island, error) {
	island, err := getIslandByID(db.DB, MyID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return island, err
}

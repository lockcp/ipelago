package database

import (
	"database/sql"
	"os"

	"github.com/ahui2016/ipelago/model"
	"github.com/ahui2016/ipelago/stmt"
	"github.com/ahui2016/ipelago/util"
	_ "github.com/mattn/go-sqlite3"
)

const MyIslandID = "My-Island-ID"

type (
	Island     = model.Island
	Message    = model.Message
	SimpleMsg  = model.SimpleMsg
	Newsletter = model.Newsletter
)

type DB struct {
	Path string
	DB   *sql.DB
}

func (db *DB) mustBegin() *sql.Tx {
	tx, err := db.DB.Begin()
	util.Panic(err)
	return tx
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

func (db *DB) CreateMyIsland(island Island) error {
	tx := db.mustBegin()
	defer tx.Rollback()

	e1 := insertIsland(tx, island)
	e2 := insertFirsstMsg(tx, island.Name)
	if err := util.WrapErrors(e1, e2); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) UpdateMyIsland(island Island) error {
	return updateIsland(db.DB, island)
}

func (db *DB) MyIsland() (Island, error) {
	island, err := getIslandByID(db.DB, MyIslandID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return island, err
}

func (db *DB) AllIslands() (islands []*Island, err error) {
	rows, err := db.DB.Query(stmt.AllIslands, MyIslandID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		island, err := scanIsland(rows)
		if err != nil {
			return nil, err
		}
		if island.Message, err = getLastMsg(db.DB, island.ID); err != nil {
			return nil, err
		}
		islands = append(islands, &island)
	}
	return islands, rows.Err()
}

func (db *DB) IslandMessages(id string) (messages []*Message, err error) {
	return getMessages(db.DB, stmt.GetIslandMessages, id)
}

func (db *DB) PostMyMsg(body string) (*Message, error) {
	if err := util.CheckStringSize(body, model.KB); err != nil {
		return nil, err
	}
	msg := model.NewMessage(body)
	if err := db.InsertMessage(msg, MyIslandID); err != nil {
		return nil, err
	}
	return msg, nil
}

func (db *DB) InsertMessage(msg *Message, id string) error {
	tx := db.mustBegin()
	defer tx.Rollback()

	if err := insertMsg(tx, msg, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) PublishNewsletter(filePath string) error {
	newsletter, err := getNewsletter(db.DB)
	if err != nil {
		return err
	}
	for length := len(newsletter); length >= model.MsgSizeLimit; {
		newsletter, err = newsletterHalf(newsletter)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, newsletter, 0644)
}

func (db *DB) InsertIsland(addr string, nl *Newsletter) error {
	tx := db.mustBegin()
	defer tx.Rollback()

	if err := nl.Trim().Check(); err != nil {
		return err
	}
	island := model.NewIsland(addr, nl)
	if err := insertIsland(tx, island); err != nil {
		return err
	}
	if err := insertMessages(tx, island.ID, nl.Messages); err != nil {
		return err
	}

	return tx.Commit()
}

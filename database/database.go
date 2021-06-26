package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ahui2016/ipelago/model"
	"github.com/ahui2016/ipelago/stmt"
	"github.com/ahui2016/ipelago/util"
	_ "github.com/mattn/go-sqlite3"
)

const MyIslandID = "My-Island-ID"

const OnePage = 10 // 每一页有多少条消息。

type (
	Island     = model.Island
	Status     = model.Status
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

func (db *DB) UpdateMyIsland(island *Island) error {
	return updateIsland(db.DB, island)
}

func (db *DB) MyIsland() (Island, error) {
	return db.GetIslandByID(MyIslandID)
}

func (db *DB) GetIslandByID(id string) (Island, error) {
	island, err := getIslandByID(db.DB, id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return island, err
}

func (db *DB) GetIslandWithoutMsg(id string) (*Island, error) {
	island, err := getIslandWithoutMsg(db.DB, id)
	return &island, err
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

// MoreIslandMessages 获取指定小岛的更多消息。
func (db *DB) MoreIslandMessages(id string, datetime int64) (messages []*Message, err error) {
	return getMessages(db.DB, stmt.GetMoreMessagesByIsland,
		id, datetime, OnePage)
}

// MoreMessages 获取全部小岛的更多消息。
func (db *DB) MoreMessages(datetime int64) (messages []*Message, err error) {
	return getMessages(db.DB, stmt.GetMoreMessages, datetime, OnePage)
}

func (db *DB) PostMyMsg(body string) (*Message, error) {
	if err := util.CheckStringSize(body, model.KB); err != nil {
		return nil, err
	}
	msg := model.NewMessage(MyIslandID, body)
	if err := db.InsertMessage(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (db *DB) InsertMessage(msg *Message) error {
	tx := db.mustBegin()
	defer tx.Rollback()

	if err := insertMsg(tx, msg); err != nil {
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

func (db *DB) InsertIsland(addr string, news *Newsletter) error {
	tx := db.mustBegin()
	defer tx.Rollback()

	if err := news.Trim().Check(); err != nil {
		return err
	}
	island := model.NewIsland(addr, news)
	if err := insertIsland(tx, island); err != nil {
		return err
	}
	if _, err := insertMessages(tx, island.ID, news.Messages); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) UpdateIsland(island *Island, news *Newsletter, oldStatus Status) (changed bool, err error) {
	tx := db.mustBegin()
	defer tx.Rollback()

	_, err = tx.Exec(stmt.UpdateIslandChecked, util.TimeNow(), island.ID)
	if err != nil {
		return
	}

	if err = news.Trim().Check(); err != nil {
		return
	}
	changed = island.UpdateFrom(news)

	// 当且只当 island 有变化（包括状态变化）时，才执行更新。
	if changed || (island.Status != oldStatus) {
		if err = updateIsland(tx, island); err != nil {
			return
		}
	}
	n, err := insertMessages(tx, island.ID, news.Messages)
	if err != nil {
		return
	}
	changed = changed || (n > 0)
	err = tx.Commit()
	return
}

func (db *DB) UpdateNote(note, id string) error {
	return db.Exec(stmt.UpdateNote, note, id)
}

func (db *DB) SetStatus(status Status, id string) error {
	return db.Exec(stmt.SetStatus, status, id)
}

func (db *DB) DeleteIsland(id string) error {
	return db.Exec(stmt.DeleteIsland, id)
}

func (db *DB) DeleteMessage(id string) error {
	n, err := getInt1(db.DB, stmt.CountMessages, MyIslandID)
	if err != nil {
		return err
	}
	if n < 2 {
		return fmt.Errorf("至少需要保留一条消息")
	}
	return db.Exec(stmt.DeleteMessage, id)
}

package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ahui2016/ipelago/stmt"
	"github.com/ahui2016/ipelago/util"
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
		&island.Email,
		&island.Avatar,
		&island.Link,
		&island.Address,
		&island.Note,
		&island.Status,
	)
	return
}

func scanMessage(row Row) (msg Message, err error) {
	err = row.Scan(
		&msg.ID,
		&msg.CTime,
		&msg.At,
		&msg.Body,
		&msg.MD,
	)
	return
}

func getIslandByID(tx TX, id string) (island Island, err error) {
	row := tx.QueryRow(stmt.GetIslandByID, id)
	if island, err = scanIsland(row); err != nil {
		return
	}
	row = tx.QueryRow(stmt.GetLastMessage, id)
	island.Message, err = scanMessage(row)
	return
}

func updateIsland(tx TX, island Island) error {
	_, err := tx.Exec(
		stmt.UpdateIsland,
		island.Name,
		island.Email,
		island.Avatar,
		island.Link,
		island.Address,
		island.Note,
		island.ID,
		island.Status,
	)
	return err
}

func insertIsland(tx TX, island Island) error {
	_, err := tx.Exec(
		stmt.InsertIsland,
		island.ID,
		island.Name,
		island.Email,
		island.Avatar,
		island.Link,
		island.Address,
		island.Note,
		island.Status,
	)
	return err
}

func insertMsg(tx TX, msg *Message, islandID string) error {
	_, e1 := tx.Exec(
		stmt.InsertMsg,
		msg.ID,
		msg.CTime,
		msg.At,
		msg.Body,
		msg.MD,
	)
	_, e2 := tx.Exec(
		stmt.InsertIslandMsg,
		islandID,
		msg.ID,
	)
	return util.WrapErrors(e1, e2)
}

// insertFirstMsg 插入每个小岛被建立时的第一条消息。
func insertFirsstMsg(tx TX, id, name string) error {
	now := time.Now()
	datetime := now.Format("2006-01-02 15:04:05")
	body := fmt.Sprintf("%s 创建于 %s", name, datetime)
	msg := &Message{
		ID:    util.RandomID(),
		CTime: now.Unix(),
		Body:  body,
	}
	return insertMsg(tx, msg, id)
}

func getMessages(tx TX, query string, args ...interface{}) (messages []*Message, err error) {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	err = rows.Err()
	return
}
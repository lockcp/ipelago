package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ahui2016/ipelago/database"
	"github.com/ahui2016/ipelago/model"
	"github.com/ahui2016/ipelago/util"
	"github.com/labstack/echo/v4"
)

// Text 用于向前端返回一个简单的文本消息。
// 为了保持一致性，总是向前端返回 JSON, 因此即使是简单的文本消息也使用 JSON.
type Text struct {
	Message string `json:"message"`
}

func errorHandler(err error, c echo.Context) {
	if e, ok := err.(*echo.HTTPError); ok {
		c.JSON(e.Code, e.Message)
	}
	c.JSON(500, Text{err.Error()})
}

func getMyIsland(c echo.Context) error {
	myIsland, err := db.MyIsland()
	if err != nil {
		return err
	}
	return c.JSON(OK, myIsland)
}

func createMyIsland(c echo.Context) error {
	name, err := getFormValue(c, "name")
	if err != nil {
		return err
	}
	email := strings.TrimSpace((c.FormValue("email")))
	avatar := strings.TrimSpace(c.FormValue("avatar"))
	link := strings.TrimSpace(c.FormValue("link"))

	island := Island{
		ID:     database.MyIslandID,
		Name:   name,
		Email:  email,
		Avatar: avatar,
		Link:   link,
	}
	return db.CreateMyIsland(island)
}

func myMessages(c echo.Context) error {
	messages, err := db.IslandMessages(database.MyIslandID)
	if err != nil {
		return err
	}
	return c.JSON(OK, messages)
}

func postMessage(c echo.Context) error {
	msgBody, err := getFormValue(c, "msg-body")
	if err != nil {
		return err
	}
	msg, err := db.PostMyMsg(msgBody)
	if err != nil {
		return err
	}
	return c.JSON(OK, msg.ID)
}

func publishNewsletter(c echo.Context) error {
	return db.PublishNewsletter(newsletterPath)
}

func followIsland(c echo.Context) (err error) {
	address := c.FormValue("address")
	done := make(chan bool, 1)

	var res *http.Response
	go func() {
		res, err = http.Get(address)
		done <- true
	}()

	var blob []byte
	select {
	case <-done:
		if err != nil { // 注意这个 err 是最外层那个 err
			return
		}
		defer res.Body.Close()

		blob, err = io.ReadAll(
			io.LimitReader(res.Body, model.MsgSizeLimit+model.KB))
		if err != nil {
			return
		}
		if err = util.CheckStringSize(string(blob), model.MsgSizeLimit); err != nil {
			return fmt.Errorf("the size exceeds the limit (15KB)")
		}
		var island model.Newsletter
		if err = json.Unmarshal(blob, &island); err != nil {
			return
		}
		if err := db.InsertIsland(address, &island); err != nil {
			return err
		}
		return c.JSON(OK, island.Name)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout")
	}
}

// getFormValue gets the c.FormValue(key), trims its spaces,
// and checks if it is empty or not.
func getFormValue(c echo.Context, key string) (string, error) {
	value := strings.TrimSpace(c.FormValue(key))
	if value == "" {
		return "", fmt.Errorf("form value [%s] is empty", key)
	}
	return value, nil
}

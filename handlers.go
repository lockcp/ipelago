package main

import (
	"fmt"
	"strings"

	"github.com/ahui2016/ipelago/database"
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
	return c.JSON(OK, myIsland)
}

func createMyIsland(c echo.Context) error {
	name, e1 := getFormValue(c, "name")
	email, e2 := getFormValue(c, "email")
	if err := util.WrapErrors(e1, e2); err != nil {
		return err
	}
	avatar := strings.TrimSpace(c.FormValue("avatar"))
	link := strings.TrimSpace(c.FormValue("link"))

	island := Island{
		ID:     database.MyIslandID,
		Name:   name,
		Email:  email,
		Avatar: avatar,
		Link:   link,
	}
	if err := db.CreateMyIsland(island); err != nil {
		return err
	}
	return restoreMyIsland()
}

func myMessages(c echo.Context) error {
	messages, err := db.IslandMessages(database.MyIslandID)
	if err != nil {
		return err
	}
	return c.JSON(OK, messages)
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

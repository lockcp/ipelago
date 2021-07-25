package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func getIslandWithoutMsg(c echo.Context) error {
	id := c.Param("id")
	island, err := db.GetIslandWithoutMsg(id)
	if err != nil {
		return err
	}
	return c.JSON(OK, island)
}

func allIslands(c echo.Context) error {
	islands, err := db.AllIslands()
	if err != nil {
		return err
	}
	return c.JSON(OK, islands)
}

func createMyIsland(c echo.Context) error {
	island, err := getFormMyIsland(c)
	if err != nil {
		return err
	}
	return db.CreateMyIsland(*island)
}

func updateMyIsland(c echo.Context) error {
	island, err := getFormMyIsland(c)
	if err != nil {
		return err
	}
	return db.UpdateMyIsland(island)
}

// moreMyMessages 获取我的小岛的更多消息。
func moreMyMessages(c echo.Context) error {
	datetime, err := getTimestamp(c)
	if err != nil {
		return err
	}
	messages, err := db.MoreIslandMessages(database.MyIslandID, datetime)
	if err != nil {
		return err
	}
	return c.JSON(OK, messages)
}

// moreIslandMessages 获取指定小岛的更多消息。
func moreIslandMessages(c echo.Context) error {
	id := c.QueryParam("id")
	datetime, err := getTimestamp(c)
	if err != nil {
		return err
	}
	messages, err := db.MoreIslandMessages(id, datetime)
	if err != nil {
		return err
	}
	return c.JSON(OK, messages)
}

// moreMessages 获取全部小岛的更多消息。
func moreMessages(c echo.Context) error {
	datetime, err := getTimestamp(c)
	if err != nil {
		return err
	}
	messages, err := db.MoreMessages(datetime)
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
	return c.JSON(OK, Text{msg.ID})
}

func publishNewsletter(c echo.Context) error {
	return db.PublishNewsletter(newsletterPath)
}

func followIsland(c echo.Context) error {
	address := c.FormValue("address")
	isDeny, e1 := db.IsDeny(address)
	isFollowed, e2 := db.IsFollowed(address)
	if err := util.WrapErrors(e1, e2); err != nil {
		return err
	}
	if isDeny {
		return fmt.Errorf("DENY")
	}
	if isFollowed {
		return fmt.Errorf("Followed")
	}
	news, err := getNews(address)
	if err != nil {
		return err
	}
	if err := db.InsertIsland(address, &news); err != nil {
		return err
	}
	return c.JSON(OK, Text{news.Name})
}

func updateIsland(c echo.Context) error {
	id, err := getFormValue(c, "id")
	if err != nil {
		return err
	}
	island, err := db.GetIslandWithoutMsg(id)
	if err != nil {
		return err
	}
	status, err := getNewsAndUpdate(island)
	if err != nil {
		return err
	}
	return c.JSON(OK, Island{Status: status})
}

func getNewsAndUpdate(island *Island) (status Status, err error) {
	// 从 island.Address 拉取消息，并根据是否超时来设置 island.Status
	oldStatus := island.Status
	news, err := getNews(island.Address)
	if util.ErrorContains(err, "timeout") {
		island.SetStatus(false)
		err = nil
	}
	if err != nil {
		return
	}
	island.SetStatus(true)

	// 尝试更新并返回 changed 表示是否真的执行了更新。
	// 如果岛名等发生了变化，或新增了消息，则 changed 为 true.
	// (注意，状态变化不影响 changed 的真假）
	changed, err := db.UpdateIsland(island, &news, oldStatus)
	if err != nil {
		return
	}
	if island.Status == model.Alive && !changed {
		island.Status = model.AliveButNoNews
	}
	return island.Status, nil
}

func getNews(address string) (news Newsletter, err error) {
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
			return
		}
		err = json.Unmarshal(blob, &news)
		return
	case <-time.After(5 * time.Second):
		return news, fmt.Errorf("timeout")
	}
}

func isCodingNet(address string) (ok bool, err error) {
	u, err := url.Parse(address)
	if err != nil {
		return
	}
	ok = strings.Contains(u.Host, "coding.net")
	return
}

func getCodingRealAddr(address string) string {
	parts := strings.Split(address, "/")
	project := parts[4]
	depot := parts[5]
	file := parts[len(parts)-1]
	return "https://" + parts[2] + "/p/" + project + "/d/" + depot + "/git/raw/master/" + file
}

func getRealAddress(address string) string {
	parts := strings.Split(address, "/")
	user := strings.Split(parts[2], ".")[0]
	return "https://" + parts[2] + "/api/user/" + user + "/project/" + parts[4] + "/shared-depot/" + parts[5] + "/git/blob/master/" + parts[len(parts)-1]
}

func updateNote(c echo.Context) error {
	id, err := getFormValue(c, "id")
	if err != nil {
		return err
	}
	note := c.FormValue("note")
	return db.UpdateNote(note, id)
}

func unfollow(c echo.Context) error {
	id, err := getFormValue(c, "id")
	if err != nil {
		return err
	}
	return db.SetStatus(model.Unfollowed, id)
}

func followAgain(c echo.Context) error {
	id, err := getFormValue(c, "id")
	if err != nil {
		return err
	}
	return db.SetStatus(model.Alive, id)
}

func deleteIsland(c echo.Context) error {
	id, err := getFormValue(c, "id")
	if err != nil {
		return err
	}
	return db.DeleteIsland(id)
}

func deleteMessage(c echo.Context) error {
	id, err := getFormValue(c, "id")
	if err != nil {
		return err
	}
	return db.DeleteMessage(id)
}

func denyIsland(c echo.Context) error {
	addr, err := getFormValue(c, "address")
	if err != nil {
		return err
	}
	return db.InsertDeny(addr)
}

func removeDeny(c echo.Context) error {
	addr, err := getFormValue(c, "address")
	if err != nil {
		return err
	}
	return db.DeleteDeny(addr)
}

func getDenyList(c echo.Context) error {
	list, err := db.GetDenyList()
	if err != nil {
		return err
	}
	return c.JSON(OK, list)
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

func getTimestamp(c echo.Context) (int64, error) {
	s := c.QueryParam("time")
	if s == "" {
		return util.TimeNow(), nil
	}
	return strconv.ParseInt(s, 10, 0)
}

func getFormMyIsland(c echo.Context) (island *Island, err error) {
	name, err := getFormValue(c, "name")
	if err != nil {
		return
	}
	email := strings.TrimSpace((c.FormValue("email")))
	link := strings.TrimSpace(c.FormValue("link"))
	avatar := strings.TrimSpace(c.FormValue("avatar"))

	if avatar != "" {
		if err = checkAvatarSize(avatar); err != nil {
			return
		}
	}

	island = &Island{
		ID:     database.MyIslandID,
		Name:   name,
		Email:  email,
		Avatar: avatar,
		Link:   link,
	}
	return
}

func checkAvatarSize(avatar string) (err error) {
	done := make(chan bool, 1)

	var res *http.Response
	go func() {
		res, err = http.Get(avatar)
		done <- true
	}()

	var blob []byte
	select {
	case <-done:
		if err != nil { // 注意这个 err 是最外层那个 err
			return
		}
		defer res.Body.Close()

		blob, err = io.ReadAll(io.LimitReader(res.Body, model.AvatarSizeLimit+model.KB))
		if err != nil {
			return
		}
		if len(blob) > model.AvatarSizeLimit {
			return fmt.Errorf("the size exceeds the limit (500KB)")
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("timeout")
	}
}

package model

import (
	"github.com/ahui2016/ipelago/util"
)

type Status string

const (
	Alive   Status = "alive"
	Timeout Status = "timeout"
	Down    Status = "down"
)

type Island struct {
	ID      string  // primary key
	Name    string  // 岛名
	Email   string  // Email
	Avatar  string  // 头像
	Link    string  // 小岛主页或岛主博客
	Address string  // 小岛地址 (JSON 文件地址)
	Note    string  // 对该小岛的备注或评价
	Message Message // 最新一条消息
	Status  Status  // 状态
}

type Message struct {
	ID   string
	Time int64
	At   string // 用于互相 @, 暂不启用
	Body string
	MD   bool // 用于 markdown, 暂不启用
}

func NewMessage(body string) *Message {
	return &Message{
		ID:   util.RandomID(),
		Time: util.TimeNow(),
		Body: body,
	}
}

type Cluster struct {
	ID   string
	Name string
}

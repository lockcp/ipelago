package model

type Island struct {
	ID      string
	Name    string // 岛名
	Avatar  string // 头像
	Email   string
	Link    string // 小岛主页或岛主博客
	Address string // 小岛地址 (JSON 文件地址)
	Note    string // 对该小岛的备注或评价
	Message string // 最新一条消息
}

type Message struct {
	ID    string
	CTime int64
	At    string
	Body  string
	MD    bool
}

type Group struct {
	ID   string
	Name string
}

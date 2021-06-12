package main

import (
	"flag"
	"net/http"

	"github.com/ahui2016/ipelago/database"
	"github.com/ahui2016/ipelago/model"
	"github.com/ahui2016/ipelago/util"
)

type (
	Island = model.Island
)

var myIsland Island

const OK = http.StatusOK

const (
	dbFileName = "ipelago.db"
	REQUIRE    = "REQUIRE.md"
)

var (
	db   = new(database.DB)
	addr = flag.String("addr", "127.0.0.1:80", "IP address of the server")
)

func init() {
	if util.PathIsNotExist(REQUIRE) {
		panic("not found: REQUIRE.md")
	}

	flag.Parse()
	util.Panic(db.Open(dbFileName))
	util.Panic(restoreMyIsland())
}

func restoreMyIsland() (err error) {
	myIsland, err = db.MyIsland()
	return
}

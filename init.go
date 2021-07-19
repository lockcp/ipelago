package main

import (
	"flag"
	"net/http"

	"github.com/ahui2016/ipelago/database"
	"github.com/ahui2016/ipelago/model"
	"github.com/ahui2016/ipelago/util"
)

const (
	OK = http.StatusOK
)

const (
	dbFileName     = "ipelago.db"
	REQUIRE        = "REQUIRE.md"
	newsletterPath = "public/newsletter.json"
)

type (
	Island     = model.Island
	Status     = model.Status
	Newsletter = model.Newsletter
)

var (
	db   = new(database.DB)
	addr = flag.String("addr", "127.0.0.1:996", "IP address of the server")
)

func init() {
	if util.PathIsNotExist(REQUIRE) {
		panic("not found: REQUIRE.md")
	}
	flag.Parse()
	util.Panic(db.Open(dbFileName))
}

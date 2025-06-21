package main

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"log/slog"
	"os"
)

var searcher *xdb.Searcher

const dbFile = "ip2region.xdb"

func init() {
	db, err := xdb.LoadContentFromFile(dbFile)
	if err != nil {
		msg := err.Error()
		cwd, err := os.Getwd()
		if err != nil {
			cwd = ""
		}
		slog.Error("can't load "+dbFile, slog.String("err", msg), slog.String("dir", cwd))
		os.Exit(1)
	}
	searcher, err = xdb.NewWithBuffer(db)
	if err != nil {
		slog.Error("can't init ip2region searcher", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

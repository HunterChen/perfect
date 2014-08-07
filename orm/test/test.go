package ormtest

import (
	"github.com/vpetrov/perfect/orm"
	"net/url"
	"os"
	"testing"
)

const (
	defaultDBName = "test1"
	defaultDBUrl  = "mongodb://localhost/" + defaultDBName
)

var DbUrl string

func NewTestDatabase(dburl string, t *testing.T) (db orm.Database, clean func()) {
	u, err := url.Parse(dburl)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	db, err = orm.NewDatabase(u, defaultDBName)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	err = db.Connect()
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	clean = func() {
		err = db.Disconnect()
		if err != nil {
			t.Fatalf("err = %v", err)
		}
	}

	return
}

func init() {
	envDBUrl := os.Getenv("DBURL")
	if len(envDBUrl) != 0 {
		DbUrl = envDBUrl
	} else {
		DbUrl = defaultDBUrl
	}
}

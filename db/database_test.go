package db

import (
	"testing"
)

func TestNewDatabase(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)
	db, err := NewDatabase(u, "")

	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if db == nil {
		t.Fatalf("db is nil, expected non-nil value")
	}
}

func TestNewDatabase_BadScheme(t *testing.T) {
	u := newTestUrl("bad_scheme://" + testHost + "/" + testDatabaseName)
	db, err := NewDatabase(u, "")

	if err == nil {
		t.Fatalf("err is nil, expected non-nil")
	}

	if db != nil {
		t.Fatalf("db is %v, expected nil", db)
	}
}

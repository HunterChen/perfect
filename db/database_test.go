package db

import (
	"testing"
)

func TestNewDatabase(t *testing.T) {
	u := newMockUrl("mongodb://localhost/test")
	db, err := NewDatabase(u, "")

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if db == nil {
		t.Errorf("db is nil, expected non-nil value")
	}
}

func TestNewDatabase_BadScheme(t *testing.T) {
	u := newMockUrl("bad_scheme://no_such_host.localdomain/")
	db, err := NewDatabase(u, "")

	if err == nil {
		t.Errorf("err is nil, expected non-nil")
	}

	if db != nil {
		t.Errorf("db is %v, expected nil", db)
	}
}

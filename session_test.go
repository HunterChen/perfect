package perfect

import (
	"testing"
    "github.com/vpetrov/perfect/orm"
)

var mock_session *Session = &Session{
    Object: orm.Object{Id:1},
	Id:            db.String("ABCD"),
	Authenticated: db.Bool(true),
	Values:        &map[string]string{"id": "ABCD", "authenticated": "1"},
}

func TestNewSession(t *testing.T) {
    session_id := "test"
	session := NewSession(session_id)

	if session.Id == nil {
        t.Fatalf("session.Id is nil, expected non-nil")
	}

    if *session.Id != session_id {
        t.Fatalf("*session.Id is '%v', expected '%v'", *session.Id, session_id)
    }

    if session.Authenticated == nil {
        t.Fatalf("session.Authenticated is nil, expected non-nil")
    }

	if *session.Authenticated {
		t.Fatalf("*session.Authenticated is %v, want %v", *session.Authenticated, false)
	}

	if session.Values == nil {
		t.Fatalf("session.Values is not allocated")
	}
}


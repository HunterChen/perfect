package perfect

import (
	"github.com/vpetrov/perfect/orm"
	"reflect"
	"testing"
)

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

func TestSession_Partial(t *testing.T) {

	type testCase struct {
		Update, Expected *Session
	}

	var (
		oid  orm.Object         = orm.Object{Id: 1}
		sid1 *string            = orm.String("1")
		sid2 *string            = orm.String("2")
		yes  *bool              = orm.Bool(true)
		no   *bool              = orm.Bool(false)
		val1 *map[string]string = &map[string]string{"1": "2"}
		val2 *map[string]string = &map[string]string{"1": "3"}
		val3 *map[string]string = &map[string]string{"2": "1"}

		partial_session_updates []testCase = []testCase{
			{Update: &Session{Object: oid}, Expected: &Session{Object: oid}},
			{Update: &Session{Object: oid, Id: sid1}, Expected: &Session{Object: oid, Id: sid1}},
			{Update: &Session{Object: oid, Id: sid2}, Expected: &Session{Object: oid, Id: sid2}},
			{Update: &Session{Object: oid, Id: sid1}, Expected: &Session{Object: oid, Id: sid1}},
			{Update: &Session{Object: oid, Id: sid1, Authenticated: no}, Expected: &Session{Object: oid, Id: sid1, Authenticated: no}},
			{Update: &Session{Object: oid, Id: sid1, Authenticated: yes}, Expected: &Session{Object: oid, Id: sid1, Authenticated: yes}},
			{Update: &Session{Object: oid, Id: sid2, Authenticated: no}, Expected: &Session{Object: oid, Id: sid2, Authenticated: no}},
			{Update: &Session{Object: oid, Id: sid2, Authenticated: yes}, Expected: &Session{Object: oid, Id: sid2, Authenticated: yes}},
			{Update: &Session{Object: oid, Id: sid2}, Expected: &Session{Object: oid, Id: sid2, Authenticated: yes}},
			{Update: &Session{Object: oid, Id: sid1}, Expected: &Session{Object: oid, Id: sid1, Authenticated: yes}},
			{Update: &Session{Object: oid, Authenticated: no}, Expected: &Session{Object: oid, Id: sid1, Authenticated: no}},
			{Update: &Session{Object: oid, Values: val1}, Expected: &Session{Object: oid, Id: sid1, Authenticated: no, Values: val1}},
			{Update: &Session{Object: oid, Values: val2}, Expected: &Session{Object: oid, Id: sid1, Authenticated: no, Values: val2}},
			{Update: &Session{Object: oid, Values: val1}, Expected: &Session{Object: oid, Id: sid1, Authenticated: no, Values: val1}},
			{Update: &Session{Object: oid, Values: val2}, Expected: &Session{Object: oid, Id: sid1, Authenticated: no, Values: val2}},
			{Update: &Session{Object: oid, Id: sid2}, Expected: &Session{Object: oid, Id: sid2, Authenticated: no, Values: val2}},
			{Update: &Session{Object: oid, Authenticated: yes}, Expected: &Session{Object: oid, Id: sid2, Authenticated: yes, Values: val2}},
			{Update: &Session{Object: oid, Values: val3}, Expected: &Session{Object: oid, Id: sid2, Authenticated: yes, Values: val3}},
			{Update: &Session{Object: oid}, Expected: &Session{Object: oid, Id: sid2, Authenticated: yes, Values: val3}},
		}
	)

	db, clean := NewTestDatabase(dbUrl, t)
	defer clean()

	//clean the collection at the start and end of this test
	db.DropCollection(partial_session_updates[0].Update)
	//defer db.DropCollection(partial_sessions[0])

	for i, test := range partial_session_updates {
		err := db.Save(test.Update)
		if err != nil {
			t.Fatalf("partial session %v: err = %v", i+1, err)
		}

		actual := &Session{
			Object: test.Update.Object,
		}

		err = db.Find(actual)
		if err != nil {
			t.Fatalf("partial session %v: err = %v", i+1, err)
		}

		if !reflect.DeepEqual(actual, test.Expected) {
			t.Fatalf("partial session %v: actual session is not exactly the same as expected session\n actual: %v\n expected: %v\n update: %v", i+1, actual, test.Expected, test.Update)
		}
	}
}

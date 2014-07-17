package db

import (
	"reflect"
	"testing"
)

type mockPerson struct {
	Object `bson:",inline,omitempty"`
	Email  string `bson:"email,omitempty"`
	Name   string `bson:"name,omitempty"`
}

func TestMongoDBQuery_One(t *testing.T) {
	db, clean := newTestMongoDB(t)
	defer clean()

	col := db.C("test2")

	err := col.Drop()
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	expected := &mockPerson{
		Email: "test@example.com",
		Name:  "John Smith",
	}

	err = col.Save(expected)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	actual := &mockPerson{
		Email: "test@example.com",
	}

	db.SetDebug(false)
	err = col.Find(actual)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("records are not equal:\nactual:  %v\nexpected: %v\n", actual, expected)
	}
}

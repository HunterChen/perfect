package orm

import (
	"reflect"
	"testing"
)

type mockPerson struct {
	Object `bson:",inline,omitempty"`
	Email  string `bson:"email,omitempty"`
	Name   string `bson:"name,omitempty"`
}

type mockUser struct {
	Object       `bson:",inline,omitempty"`
	Id           *string   `bson:"id,omitempty"`
	Email        *string   `bson:"email,omitempty"`
	Name         *string   `bson:"name,omitempty"`
	Organization *string   `bson:"org,omitempty"`
	Age          *int      `bson:"age,omitempty"`
	Address      *string   `bson:"address,omitempty"`
	Tags         *[]string `bson:"tags,omitempty"`
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

	err = col.Find(actual)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("records are not equal:\nactual:  %v\nexpected: %v\n", actual, expected)
	}
}

func TestMongoDBQUery_Select(t *testing.T) {
	var err error

	db, clean := newTestMongoDB(t)
	defer clean()

	original := &mockUser{
		Id:           String("abc123"),
		Email:        String("user@example.com"),
		Name:         String("John D'Oh"),
		Organization: String("NSA"),
		Age:          Int(32),
		Address:      String("123 Golang Way, Mars City, 0000001, Mars, Planet Mars, Solar System"),
		Tags:         &[]string{"user", "admin", "mars", "NSA"},
	}

	err = db.Save(original)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	//cleanup
	defer func() {
		err = db.Remove(original)
		if err != nil {
			t.Fatalf("err = %v")
		}
	}()

	expected := &mockUser{
		Object: original.Object,
		Id:     original.Id,
		Email:  original.Email,
		Tags:   original.Tags,
	}

	actual := &mockUser{
		Object: original.Object,
	}

	err = db.Query(actual).Select("id", "email", "tags").One(actual)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("records are not equal:\nactual: %v\nexpected: %v\n", actual, expected)
	}
}

func TestMongoDBQUery_Exclude(t *testing.T) {
	var err error

	db, clean := newTestMongoDB(t)
	defer clean()

	original := &mockUser{
		Id:           String("abc123"),
		Email:        String("user@example.com"),
		Name:         String("John D'Oh"),
		Organization: String("NSA"),
		Age:          Int(32),
		Address:      String("123 Golang Way, Mars City, 0000001, Mars, Planet Mars, Solar System"),
		Tags:         &[]string{"user", "admin", "mars", "NSA"},
	}

	err = db.Save(original)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	//cleanup
	defer func() {
		err = db.Remove(original)
		if err != nil {
			t.Fatalf("err = %v")
		}
	}()

	expected := &mockUser{
		Object: original.Object,
		Id:     original.Id,
		Email:  original.Email,
		Tags:   original.Tags,
	}

	actual := &mockUser{
		Object: original.Object,
	}

	err = db.Query(actual).Exclude("name", "org", "age", "address").One(actual)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("records are not equal:\nactual: %#v\nexpected: %#v\n", actual, expected)
	}
}

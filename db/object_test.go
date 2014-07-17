package db

import (
	"testing"
)

func TestGetEmptyId(t *testing.T) {
	o := Object{}
	id := o.GetDbId()
	if id != nil {
		t.Fatalf("expected an empty object's ID to be nil, got %v instead", id)
	}
}

func TestGetSetId(t *testing.T) {
	o := Object{}
	id := "my_id"

	o.SetDbId(id)

	id2 := o.GetDbId()

	if id2 != id {
		t.Fatalf("object id is %v, expected %v", id2, id)
	}
}

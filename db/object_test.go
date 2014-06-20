package db

import (
	"testing"
)

func TestGetEmptyId(t *testing.T) {
	o := Object{}
	id := o.GetDbId()
	if id != nil {
		t.Errorf("expected an empty object's ID to be nil, got %v instead", id)
	}
}

func TestGetSetId(t *testing.T) {
	o := Object{}
	id := "my_id"

	o.SetDbId(id)

	id2 := o.GetDbId()

	if id2 != id {
		t.Errorf("object id is %v, expected %v", id2, id)
	}
}

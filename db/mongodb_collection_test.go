package db

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"testing"
)

type mockRecord struct {
	Object `bson:",inline,omitempty" json:"-"`
	S      string `S`
	I      int    `I`
	/*
		I8     int8
		I16    int16
		I32    int32
		I64    int64
		BOOL   bool
		B      byte
		F32    float32
		F64    float64
		R      rune
		UI     uint
		UI8    uint8
		UI16   uint16
		UI32   uint32
		UI64   uint64
	*/
}

func printJSON(o interface{}) {
	data, err := json.Marshal(o)

	log.Println("[json]", string(data), "err=", err)
}

func cleanMongoDBCollection(c Collection) {
	c.(*MongoDBCollection).RemoveAll(nil)
}

func TestMongoDBCollection_Name(t *testing.T) {
	db, clean := newRealMongoDB(t)
	defer clean()

	cname := "test"
	col := db.C(cname)
	if col == nil {
		t.Errorf("col is nil, expected non- nil")
	}

	actual_name := col.Name()

	if actual_name != cname {
		t.Errorf("Collection name is '%v', expected '%v'", actual_name, cname)
	}
}

func TestMongoDBCollection_Drop(t *testing.T) {
	db, clean := newRealMongoDB(t)
	defer clean()

	col := db.C("test")
	err := col.Drop()
	if err != nil {
		t.Errorf("Drop returned '%v', expected nil", err)
	}
}

func TestMongoDBCollection_Save(t *testing.T) {
	var err error

	db, clean := newRealMongoDB(t)
	defer clean()

	col := db.C("test2")
	//defer cleanMongoDBCollection(col)

	original := &mockRecord{
		S: "test",
		I: 42,
	}

	printJSON(original)

	err = col.Save(original)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	id := original.GetDbId()
	if id == nil {
		t.Errorf("record id is nil, expected non-nil")
	}

	bson_id, ok := id.(bson.ObjectId)
	if !ok {
		t.Errorf("record id is %v, expected bson.ObjectId", id)
	}

	if !bson_id.Valid() {
		t.Errorf("record id is %v, expected a valid bson.ObjectId", id)
	}

	//read the entry back into a map by using mgo directly
	mc := col.(*MongoDBCollection)
	actual := map[string][]byte{}

	err = mc.FindId(id).One(&actual)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	log.Println("%#v", actual)

	//check that S is a string
	if string(actual["S"]) != original.S {
		t.Errorf("record.S is '%v', expected '%v'", string(actual["S"]), original.S)
	}

}

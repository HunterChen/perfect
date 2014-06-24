package db

import (
	"encoding/json"
	"github.com/vpetrov/perfect/pointer"
	"labix.org/v2/mgo/bson"
	"log"
	"math"
	"testing"
)

type mockRecord struct {
	Object `bson:",inline,omitempty" json:"-"`
	S      string  `S`
	I      int     `I`
	I8     int8    `I8`
	I16    int16   `I16`
	I32    int32   `I32`
	I64    int64   `I64`
	BOOL   bool    `BOOL`
	B      byte    `B`
	F32    float32 `F32`
	F64    float64 `F64`
	R      rune    `R`
	UI     uint    `UI`
	UI8    uint8   `UI8`
	UI16   uint16  `UI16`
	UI32   uint32  `UI32`
	//UI64   uint64 //bson has no UINT64 type
	M map[string]string `M`
	A []string          `A`

	PS    *string  `PS`
	PI    *int     `PI`
	PI8   *int8    `PI8`
	PI16  *int16   `PI16`
	PI32  *int32   `PI32`
	PI64  *int64   `PI64`
	PBOOL *bool    `PBOOL`
	PB    *byte    `PB`
	PF32  *float32 `PF32`
	PF64  *float64 `PF64`
	PR    *rune    `PR`
	PUI   *uint    `PUI`
	PUI8  *uint8   `PUI8`
	PUI16 *uint16  `PUI16`
	PUI32 *uint32  `PUI32`
	//UI64   uint64 //bson has no UINT64 type
	PM *map[string]string `PM`
	PA *[]string          `PA`
}

func printJSON(o interface{}) {
	data, err := json.Marshal(o)

	log.Println("[json]", string(data), "err=", err)
}

func cleanMongoDBCollection(c Collection) {
	c.(*MongoDBCollection).RemoveAll(nil)
}

func compareRecords(actual, expected *mockRecord, t *testing.T) {
	if actual.S != expected.S {
		t.Errorf("record.S is '%v', expected '%v'", actual.S, expected.S)
	}

	if actual.I != expected.I {
		t.Errorf("record.I is %v, expected %v", actual.I, expected.I)
	}

	if actual.I8 != expected.I8 {
		t.Errorf("record.I8 is %v, expected %v", actual.I8, expected.I8)
	}

	if actual.I16 != expected.I16 {
		t.Errorf("record.I16 is %v, expected %v", actual.I16, expected.I16)
	}

	if actual.I32 != expected.I32 {
		t.Errorf("record.I32 is %v, expected %v", actual.I32, expected.I32)
	}

	if actual.I64 != expected.I64 {
		t.Errorf("record.I64 is %v, expected %v", actual.I64, expected.I64)
	}

	if actual.BOOL != expected.BOOL {
		t.Errorf("record.BOOL is %v, expected %v", actual.BOOL, expected.BOOL)
	}

	if actual.F32 != expected.F32 {
		t.Errorf("record.F32 is %v, expected %v", actual.F32, expected.F32)
	}

	if actual.F64 != expected.F64 {
		t.Errorf("record.F64 is %v, expected %v", actual.F64, expected.F64)
	}

	if actual.R != expected.R {
		t.Errorf("record.R is '%v', expected '%v'", actual.R, expected.R)
	}

	if actual.UI != expected.UI {
		t.Errorf("record.UI is %v, expected %v", actual.UI, expected.UI)
	}

	if actual.UI8 != expected.UI8 {
		t.Errorf("record.UI8 is %v, expected %v", actual.UI8, expected.UI8)
	}

	if actual.UI16 != expected.UI16 {
		t.Errorf("record.UI16 is %v, expected %v", actual.UI16, expected.UI16)
	}

	if actual.UI32 != expected.UI32 {
		t.Errorf("record.UI32 is %v, expected %v", actual.UI32, expected.UI32)
	}

	//iterate over the map
	for k, expected_v := range expected.M {
		actual_v, ok := actual.M[k]
		if !ok {
			t.Errorf("record.M[\"%v\"] doesn't exist, expected it to exist.", k)
		}

		if actual_v != expected_v {
			t.Errorf("record.M[\"%v\"] is '%v', expected '%v'", k, actual_v, expected_v)
		}
	}

	//verify number of elements
	nexpectedA := len(expected.A)
	nactualA := len(actual.A)

	if nactualA != nexpectedA {
		t.Errorf("len(record.A) is %v, expected %v\nrecord.A: %v\nexpected.A: %v", nactualA, nexpectedA, actual.A, expected.A)
	}

	//verify each element
	for i := 0; i < nexpectedA; i++ {
		if actual.A[i] != expected.A[i] {
			t.Errorf("record.A[%v] is %v, expected %v", i, actual.A[i], expected.A[i])
		}
	}

	/* Verify pointer types */

	//returns "nil" or "non-nil" depending on whether 'a' is nil
	nilOrNot := func(a interface{}) string {
		if a == nil {
			return "nil"
		}

		return "non-nil"
	}

	if actual.PS != expected.PS {
		if actual.PS == nil || expected.PS == nil {
			t.Errorf("record.PS is %v, expected %v", nilOrNot(actual.PS), nilOrNot(expected.PS))
		} else if *actual.PS != *expected.PS {
			t.Errorf("*record.S is '%v', expected '%v'", *actual.PS, *expected.PS)
		}
	}

	if actual.PI != expected.PI {
		if actual.PI == nil || expected.PI == nil {
			t.Errorf("record.PI is %v, expected %v", nilOrNot(actual.PI), nilOrNot(expected.PI))
		} else if *actual.PI != *expected.PI {
			t.Errorf("*record.PI is %v, expected %v", *actual.PI, *expected.PI)
		}
	}

	if actual.PI8 != expected.PI8 {
		if actual.PI8 == nil || expected.PI8 == nil {
			t.Errorf("record.PI8 is %v, expected %v", nilOrNot(actual.PI8), nilOrNot(expected.PI8))
		} else if *actual.PI8 != *expected.PI8 {
			t.Errorf("*record.PI8 is %v, expected %v", *actual.PI8, *expected.PI8)
		}
	}

	if actual.PI16 != expected.PI16 {
		if actual.PI16 == nil || expected.PI16 == nil {
			t.Errorf("record.PI16 is %v, expected %v", nilOrNot(actual.PI16), nilOrNot(expected.PI16))
		} else if *actual.PI16 != *expected.PI16 {
			t.Errorf("record.PI16 is %v, expected %v", *actual.PI16, *expected.PI16)
		}
	}

	if actual.PI32 != expected.PI32 {
		if actual.PI32 == nil || expected.PI32 == nil {
			t.Errorf("record.PI32 is %v, expected %v", nilOrNot(actual.PI32), nilOrNot(actual.PI32))
		} else if *actual.PI32 != *expected.PI32 {
			t.Errorf("record.PI32 is %v, expected %v", *actual.PI32, *expected.PI32)
		}
	}

	if actual.PI64 != expected.PI64 {
		if actual.PI64 == nil || expected.PI64 == nil {
			t.Errorf("record.PI64 is %v, expected %v", nilOrNot(actual.PI64), nilOrNot(actual.PI64))
		} else if *actual.PI64 != *expected.PI64 {
			t.Errorf("record.PI64 is %v, expected %v", *actual.PI64, *expected.PI64)
		}
	}

	if actual.PBOOL != expected.PBOOL {
		if actual.PBOOL == nil || expected.PBOOL == nil {
			t.Errorf("record.PBOOL is %v, expected %v", nilOrNot(actual.PBOOL), nilOrNot(actual.PBOOL))
		} else if *actual.PBOOL != *expected.PBOOL {
			t.Errorf("record.PBOOL is %v, expected %v", *actual.PBOOL, *expected.PBOOL)
		}
	}

	if actual.PF32 != expected.PF32 {
		if actual.PF32 == nil || expected.PF32 == nil {
			t.Errorf("record.PF32 is %v, expected %v", nilOrNot(actual.PF32), nilOrNot(actual.PF32))
		} else if *actual.PF32 != *expected.PF32 {
			t.Errorf("record.PF32 is %v, expected %v", *actual.PF32, *expected.PF32)
		}
	}

	if actual.PF64 != expected.PF64 {
		if actual.PF64 == nil || expected.PF64 == nil {
			t.Errorf("record.PF64 is %v, expected %v", nilOrNot(actual.PF64), nilOrNot(actual.PF64))
		} else if *actual.PF64 != *expected.PF64 {
			t.Errorf("record.PF64 is %v, expected %v", *actual.PF64, *expected.PF64)
		}
	}

	if actual.PR != expected.PR {
		if actual.PR == nil || expected.PR == nil {
			t.Errorf("record.PR is %v, expected %v", nilOrNot(actual.PR), nilOrNot(actual.PR))
		} else if *actual.PR != *expected.PR {
			t.Errorf("record.PR is '%v', expected '%v'", *actual.PR, *expected.PR)
		}
	}

	if actual.PUI != expected.PUI {
		if actual.PUI == nil || expected.PUI == nil {
			t.Errorf("record.PUI is %v, expected %v", nilOrNot(actual.PUI), nilOrNot(actual.PUI))
		} else if *actual.PUI != *expected.PUI {
			t.Errorf("record.PUI is %v, expected %v", *actual.PUI, *expected.PUI)
		}
	}

	if actual.PUI8 != expected.PUI8 {
		if actual.PUI8 == nil || expected.PUI8 == nil {
			t.Errorf("record.PUI8 is %v, expected %v", nilOrNot(actual.PUI8), nilOrNot(actual.PUI8))
		} else if *actual.PUI8 != *expected.PUI8 {
			t.Errorf("record.PUI8 is %v, expected %v", *actual.PUI8, *expected.PUI8)
		}
	}

	if actual.PUI16 != expected.PUI16 {
		if actual.PUI16 == nil || expected.PUI16 == nil {
			t.Errorf("record.PUI16 is %v, expected %v", nilOrNot(actual.PUI16), nilOrNot(actual.PUI16))
		} else if *actual.PUI16 != *expected.PUI16 {
			t.Errorf("record.PUI16 is %v, expected %v", *actual.PUI16, *expected.PUI16)
		}
	}

	if actual.PUI32 != expected.PUI32 {
		if actual.PUI32 == nil || expected.PUI32 == nil {
			t.Errorf("record.PUI32 is %v, expected %v", nilOrNot(actual.PUI32), nilOrNot(actual.PUI32))
		} else if *actual.PUI32 != *expected.PUI32 {
			t.Errorf("record.PUI32 is %v, expected %v", *actual.PUI32, *expected.PUI32)
		}
	}

	if actual.PM != expected.PM {
		if actual.PM == nil || expected.PM == nil {
			t.Errorf("record.PM is %v, expected %v", nilOrNot(actual.PM), nilOrNot(actual.PM))
		} else {
			//iterate over the map
			for k, expected_v := range *expected.PM {
				actual_v, ok := (*actual.PM)[k]
				if !ok {
					t.Errorf("record.PM[\"%v\"] doesn't exist, expected it to exist.", k)
				}

				if actual_v != expected_v {
					t.Errorf("record.PM[\"%v\"] is '%v', expected '%v'", k, actual_v, expected_v)
				}
			}
		}
	}

	if actual.PM != expected.PM {
		if actual.PM == nil || expected.PM == nil {
			t.Errorf("record.PM is %v, expected %v", nilOrNot(actual.PM), nilOrNot(actual.PM))
		} else {
			//verify number of elements
			nexpectedA := len(*expected.PA)
			nactualA := len(*actual.PA)

			if nactualA != nexpectedA {
				t.Errorf("len(record.PA) is %v, expected %v\nrecord.PA: %v\nexpected.PA: %v", nactualA, nexpectedA, *actual.PA, *expected.PA)
			}

			//verify each element
			for i := 0; i < nexpectedA; i++ {
				if (*actual.PA)[i] != (*expected.PA)[i] {
					t.Errorf("record.PA[%v] is %v, expected %v", i, (*actual.PA)[i], (*expected.PA)[i])
				}
			}
		}
	}
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

	expected := &mockRecord{
		S:    "test",
		I:    math.MaxInt32,
		I8:   math.MaxInt8,
		I16:  math.MaxInt16,
		I32:  math.MaxInt32,
		I64:  math.MaxInt64,
		BOOL: true,
		B:    255,
		F32:  math.MaxFloat32,
		F64:  math.MaxFloat64,
		R:    math.MaxInt32,
		UI:   math.MaxUint32,
		UI8:  math.MaxUint8,
		UI16: math.MaxUint16,
		UI32: math.MaxUint32,
		//UI64: math.MaxUint64, //there is no UINT64 in BSON
		M: map[string]string{
			"key 1": "value 1",
			"key 2": "value 2",
		},
		A: []string{
			"value 1", "value 2", "value 3",
		},
		PS:    pointer.String("test"),
		PI:    pointer.Int(math.MaxInt32),
		PI8:   pointer.Int8(math.MaxInt8),
		PI16:  pointer.Int16(math.MaxInt16),
		PI32:  pointer.Int32(math.MaxInt32),
		PI64:  pointer.Int64(math.MaxInt64),
		PBOOL: pointer.Bool(true),
		PB:    pointer.Byte(255),
		PF32:  pointer.Float32(math.MaxFloat32),
		PF64:  pointer.Float64(math.MaxFloat64),
		PR:    pointer.Rune(math.MaxInt32),
		PUI:   pointer.Uint(math.MaxUint32),
		PUI8:  pointer.Uint8(math.MaxUint8),
		PUI16: pointer.Uint16(math.MaxUint16),
		PUI32: pointer.Uint32(math.MaxUint32),
		//UI64: math.MaxUint64, //there is no UINT64 in BSON
		PM: &map[string]string{
			"key 1": "value 1",
			"key 2": "value 2",
		},
		PA: &[]string{
			"value 1", "value 2", "value 3",
		},
	}

	printJSON(expected)

	err = col.Save(expected)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	id := expected.GetDbId()
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
	actual := &mockRecord{}

	err = mc.FindId(id).One(&actual)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	log.Println("%#v", actual)

	compareRecords(actual, expected, t)
}

func TestA1(t *testing.T) {
	a := &mockRecord{}
	b := &mockRecord{}
	compareRecords(a, b, t)
}

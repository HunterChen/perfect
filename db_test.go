package perfect

import (
        "testing"
        "net/url"
)

func TestDBO_DbId(t *testing.T) {
    DB_ID := 123
    DB_COL := "test"

    dbo := DBO{
            DBID: DB_ID,
            Collection : DB_COL,
         }

    if (dbo.DbId() != DB_ID) {
		t.Errorf("dbo.DbId() returns %v, want %v", dbo.DbId(), DB_ID)
    }
}

func TestDBO_SetDbId(t *testing.T) {
    DB_ID := 123
    NEW_DB_ID := 321
    DB_COL := "test"

    dbo := DBO{
            DBID: DB_ID,
            Collection : DB_COL,
         }

    dbo.SetDbId(NEW_DB_ID)

    if (dbo.DbId() != NEW_DB_ID) {
		t.Errorf("dbo.SetDbId() was supposed to the DBID value to %v, instead DbId() returned %v", NEW_DB_ID, dbo.DbId())
    }
}

func TestDBO_DbCollection(t *testing.T) {
    DB_COL := "test"

    dbo := DBO{
            Collection : DB_COL,
         }

    if (dbo.DbCollection() != DB_COL) {
		t.Errorf("dbo.DbCollection() returned %v, expected %v", dbo.DbCollection, DB_COL)
    }
}

func TestDBO_SetDbCollection(t *testing.T) {
    DB_COL := "test"
    NEW_DB_COL := "tset"

    dbo := DBO{
            Collection : DB_COL,
         }

    dbo.SetDbCollection(NEW_DB_COL)

    if (dbo.DbCollection() != NEW_DB_COL) {
		t.Errorf("dbo.DbSetCollection() was supposed to set the Collection value %v, instead DbCollection() returned %v", NEW_DB_COL, dbo.DbCollection())
    }
}

func TestNewDatabase_MongoDB(t *testing.T) {
    DB_NAME := "test"
    u, err := url.Parse("mongodb://localhost:27017/test")
    if err != nil {
		t.Errorf("err = %v", err)
    }

    var db Database = NewDatabase(u, DB_NAME)

    if db == nil {
		t.Errorf("NewDatabase returned nil, want non-nil value")
    }
}

func TestNewDatabase_NIL(t *testing.T) {
    DB_NAME := "test"
    u, err := url.Parse("test:///")
    if err != nil {
		t.Errorf("err = %v", err)
    }

    var db Database = NewDatabase(u, DB_NAME)

    if db != nil {
		t.Errorf("NewDatabase returned %v, expected nil", db)
    }
}

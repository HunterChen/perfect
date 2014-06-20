package db

import (
	"net/url"
	"testing"
)

/**********
* HELPERS *
**********/

func newMockUrl(u string) (result *url.URL) {
	result, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	return
}

func mockMongoDBUrl(dbname string) *url.URL {
	return newMockUrl("mongodb://localhost/" + dbname)
}

func realMongoDBUrl(dbname string) *url.URL {
	return newMockUrl("mongodb://localhost:27017/" + dbname)
}

func badMongoDBUrl(dbname string) *url.URL {
	return newMockUrl("mongodb://no_such_host.localdomain/" + dbname)
}

func newMockMongoDB(t *testing.T) *MongoDB {
	dbname := "test"
	u := mockMongoDBUrl(dbname)
	db := NewMongoDB(u, dbname)
	if db == nil {
		t.Errorf("db is nil, expected non-nil")
	}

	return db
}

//returns a new connection to MongoDB and a clean-up function which the caller
//should use with defer.
func newRealMongoDB(t *testing.T) (*MongoDB, func()) {
	u := realMongoDBUrl("test1")
	db := NewMongoDB(u, "")
	err := db.Connect()
	if err != nil {
		t.Errorf("err = %v", err)
	}

	cleanup := func() {
		err := db.Disconnect()
		if err != nil {
			t.Errorf("err = %v", err)
		}
	}

	return db, cleanup
}

/********
* TESTS *
********/

func TestNewMongoDB(t *testing.T) {
	u := mockMongoDBUrl("")
	dbname := "test"

	db := NewMongoDB(u, dbname)
	if db == nil {
		t.Errorf("db is nil, expected a non-nil value")
	}
}

func TestMongoDB_NameFromArgument(t *testing.T) {
	dbname := "test"
	u := mockMongoDBUrl("")
	db := NewMongoDB(u, dbname)

	actual_dbname := db.Name()
	if actual_dbname != dbname {
		t.Errorf("database name is %v, expected %v", actual_dbname, dbname)
	}
}

func TestMongoDB_NameFromURL(t *testing.T) {
	dbname := "test_db"
	u := mockMongoDBUrl(dbname)

	db := NewMongoDB(u, "")
	if db == nil {
		t.Errorf("db is nil, expected a non-nil value")
	}

	actual_dbname := db.Name()
	if actual_dbname != dbname {
		t.Errorf("database name is %v, expected %v", actual_dbname, dbname)
	}
}

func TestMongoDB_NameFromArgumentAndURL(t *testing.T) {
	dbname := "test_db"
	u := mockMongoDBUrl(dbname)

	db := NewMongoDB(u, dbname)
	actual_dbname := db.Name()

	if actual_dbname != dbname {
		t.Errorf("database name is %v, expected %v", actual_dbname, dbname)
	}
}

func TestMongoDB_Url(t *testing.T) {
	u := mockMongoDBUrl("test")

	db := NewMongoDB(u, "")
	if db == nil {
		t.Errorf("db is nil, expected a non-nil value")
	}

	actual_url := db.URL()

	if actual_url != u {
		t.Errorf("database url is %v, expected %v", actual_url, u)
	}
}

func TestMongoDB_ConnectAndDisconnect(t *testing.T) {
	dbname := "test"
	u := realMongoDBUrl(dbname)
	db := NewMongoDB(u, "")

	err := db.Connect()
	if err != nil {
		t.Errorf("Connect() returned an error: %v", err)
	}

	err = db.Disconnect()
	if err != nil {
		t.Errorf("Disconnect() returned an error: %v", err)
	}
}

//tests a connection attempt to an invalid mongodb instance.
//may take 10 sec to complete
func TestMongoDB_ConnectToUnreachableServer(t *testing.T) {
	u := newMockUrl("mongodb://127.0.0.1:0/")
	db := NewMongoDB(u, "")

	err := db.Connect()
	if err == nil {
		t.Errorf("Connect(): expected an error")
	}

	err = db.Disconnect()
	if err != nil {
		t.Errorf("Disconnect(): err = %v", err)
	}
}

func TestMongoDB_SystemInformationOffline(t *testing.T) {
	u := mockMongoDBUrl("")
	db := NewMongoDB(u, "")

	info := db.SystemInformation()
	if len(info) != 0 {
		t.Errorf("system information is '%v', expected '%v'", info, "")
	}
}

func TestMongoDB_SystemInformationOnline(t *testing.T) {
	db, clean := newRealMongoDB(t)
	defer clean()

	info := db.SystemInformation()
	if len(info) == 0 {
		t.Errorf("system information is empty, expected a non-empty string")
	}
}

func TestMongoDB_VersionOffline(t *testing.T) {
	db := newMockMongoDB(t)

	version := db.Version()
	if len(version) != 0 {
		t.Errorf("version is '%v', expected '%v'", version, "")
	}
}

func TestMongoDB_VersionOnline(t *testing.T) {
	db, clean := newRealMongoDB(t)
	defer clean()

	version := db.Version()
	if len(version) == 0 {
		t.Errorf("version is empty, expected a non-empty string")
	}
}

func TestMongoDB_C_Offline(t *testing.T) {
	db := newMockMongoDB(t)
	c := db.C("test")

	if c != nil {
		t.Errorf("collection is %v, expected nil", c)
	}
}

func TestMongoDB_C_Online(t *testing.T) {
	db, clean := newRealMongoDB(t)
	defer clean()

	c := db.C("test")

	if c == nil {
		t.Errorf("collection is nil, expected non-nil")
	}
}

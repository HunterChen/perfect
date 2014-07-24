package orm

import (
	"labix.org/v2/mgo"
	"net/url"
	"testing"
)

/**********
* HELPERS *
**********/

const (
	testHost           = "localhost:27017"
	testHostUrl        = "mongodb://" + testHost + "/"
	testDatabaseName   = "test1"
	testDatabaseUrl    = testHostUrl + testDatabaseName
	testCollectionName = "test"
	testUsername       = "test"
	testPassword       = "test"
	testAuthUrl        = "mongodb://" + testUsername + ":" + testPassword + "@" + testHost + "/" + testDatabaseName

	//'invalid'
	testInvalidAuthUrl     = "mongodb://" + testUsername + ":" + testPassword + "invalid" + "@" + testHost + "/" + testDatabaseName
	testInvalidHost        = "no_such_host.localdomain"
	testInvalidUrl         = "mongodb://" + testInvalidHost + "/"
	testInvalidDatabaseUrl = testInvalidUrl + testDatabaseName

	testUnreachableMongoDBHost    = "localhost:0"
	testUnreachableMongoDBHostUrl = "mongodb://" + testUnreachableMongoDBHost + "/"
)

func newTestUrl(u string) (result *url.URL) {
	result, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	return
}

//returns a new connection to MongoDB and a clean-up function which the caller
//should use with defer.
func newTestMongoDB(t *testing.T) (*MongoDB, func()) {
	u := newTestUrl(testDatabaseUrl)
	return newTestMongoDBWithURL(u, t)
}

func newTestMongoDBWithURL(u *url.URL, t *testing.T) (*MongoDB, func()) {
	db := NewMongoDB(u, "")
	err := db.Connect()
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	cleanup := func() {
		err := db.Disconnect()
		if err != nil {
			t.Fatalf("err = %v", err)
		}
	}

	return db, cleanup
}

func setupAuth(db *MongoDB, username, password string, t *testing.T) func() {
	user := &mgo.User{
		Username: username,
		Password: password,
		Roles:    []mgo.Role{mgo.RoleReadWrite},
	}

	err := db.Database.UpsertUser(user)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	return func() {
		err := db.Database.RemoveUser(username)
		if err != nil {
			t.Fatalf("err = %v")
		}
	}
}

/********
* TESTS *
********/

func TestNewMongoDB(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)

	db := NewMongoDB(u, testDatabaseName)
	if db == nil {
		t.Fatalf("db is nil, expected a non-nil value")
	}
}

func TestMongoDB_NameFromArgument(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)
	db := NewMongoDB(u, testDatabaseName)

	actual_dbname := db.Name()
	if actual_dbname != testDatabaseName {
		t.Fatalf("database name is %v, expected %v", actual_dbname, testDatabaseName)
	}
}

func TestMongoDB_NameFromURL(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)

	db := NewMongoDB(u, "")
	if db == nil {
		t.Fatalf("db is nil, expected a non-nil value")
	}

	actual_dbname := db.Name()
	if actual_dbname != testDatabaseName {
		t.Fatalf("database name is %v, expected %v", actual_dbname, testDatabaseName)
	}
}

func TestMongoDB_NameFromArgumentAndURL(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)

	db := NewMongoDB(u, testDatabaseName)
	actual_dbname := db.Name()

	if actual_dbname != testDatabaseName {
		t.Fatalf("database name is %v, expected %v", actual_dbname, testDatabaseName)
	}
}

func TestMongoDB_Url(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)

	db := NewMongoDB(u, "")
	if db == nil {
		t.Fatalf("db is nil, expected a non-nil value")
	}

	actual_url := db.URL()

	if actual_url != u {
		t.Fatalf("database url is %v, expected %v", actual_url, u)
	}
}

func TestMongoDB_ConnectAndDisconnect(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)
	db := NewMongoDB(u, "")

	err := db.Connect()
	if err != nil {
		t.Fatalf("Connect() returned an error: %v", err)
	}

	err = db.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect() returned an error: %v", err)
	}
}

//tests a connection attempt to an invalid mongodb instance.
//may take 10 sec to complete
func TestMongoDB_ConnectToUnreachableServer(t *testing.T) {
	u := newTestUrl(testUnreachableMongoDBHostUrl)
	db := NewMongoDB(u, "")

	err := db.Connect()
	if err == nil {
		t.Fatalf("Connect(): expected an error")
	}

	err = db.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect(): err = %v", err)
	}
}

func TestMongoDB_ConnectWithUser(t *testing.T) {

	var (
		err error
	)

	db, clean := newTestMongoDB(t)
	defer clean()

	clean = setupAuth(db, testUsername, testPassword, t)
	defer clean()

	dburl := newTestUrl(testAuthUrl)

	//create another connection to MongoDB
	db = NewMongoDB(dburl, "")
	err = db.Connect()
	if err != nil {
		t.Fatalf("Connect(): err is '%v', expected nil", err)
	}

	err = db.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect(): err is '%v', expected nil", err)
	}
}

func TestMongoDB_ConnectWithInvalidUser(t *testing.T) {
	var (
		err error
	)

	db, clean := newTestMongoDB(t)
	defer clean()

	clean = setupAuth(db, testUsername, testPassword, t)
	defer clean()

	dburl := newTestUrl(testInvalidAuthUrl)

	//create another connection to MongoDB
	db = NewMongoDB(dburl, "")
	err = db.Connect()
	if err == nil {
		t.Fatalf("Connect(): err is nil, expected non-nil")
	}

	err = db.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect(): err is '%v', expected nil", err)
	}
}

func TestMongoDB_SystemInformationOffline(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)
	db := NewMongoDB(u, "")

	info := db.SystemInformation()
	if len(info) != 0 {
		t.Fatalf("system information is '%v', expected '%v'", info, "")
	}
}

func TestMongoDB_SystemInformationOnline(t *testing.T) {
	db, clean := newTestMongoDB(t)
	defer clean()

	info := db.SystemInformation()
	if len(info) == 0 {
		t.Fatalf("system information is empty, expected a non-empty string")
	}
}

func TestMongoDB_VersionOffline(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)
	db := NewMongoDB(u, "")

	version := db.Version()
	if len(version) != 0 {
		t.Fatalf("version is '%v', expected '%v'", version, "")
	}
}

func TestMongoDB_VersionOnline(t *testing.T) {
	db, clean := newTestMongoDB(t)
	defer clean()

	version := db.Version()
	if len(version) == 0 {
		t.Fatalf("version is empty, expected a non-empty string")
	}
}

func TestMongoDB_C_Offline(t *testing.T) {
	u := newTestUrl(testDatabaseUrl)
	db := NewMongoDB(u, "")
	c := db.C("test")

	if c != nil {
		t.Fatalf("collection is %v, expected nil", c)
	}
}

func TestMongoDB_C_Online(t *testing.T) {
	db, clean := newTestMongoDB(t)
	defer clean()

	c := db.C("test")

	if c == nil {
		t.Fatalf("collection is nil, expected non-nil")
	}
}

func TestMongoDB_SetDebug(t *testing.T) {
	db, clean := newTestMongoDB(t)
	defer clean()

	db.SetDebug(true)
	db.SetDebug(false)
}

func TestMongoDB_UniqueId(t *testing.T) {
	db, clean := newTestMongoDB(t)
	defer clean()

	ids := []string{}

	//a helper function that checks whether a string already exists in the slice
	id_exists := func(dest []string, id string) bool {
		for _, v := range dest {
			if v == id {
				return true
			}
		}

		return false
	}

	//generate 100 IDs
	for i := 0; i < 100; i++ {
		new_id := db.UniqueId()
		if id_exists(ids, new_id) {
			t.Fatalf("duplicate id generated: %v\ncurrent ids: %v", new_id, ids)
		}

		ids = append(ids, new_id)
	}
}

func TestMongoDB_GetCollectionName(t *testing.T) {
	type testObject struct {
		Object
	}

	o := &testObject{}
	expected_col := "testobjects"

	db, clean := newTestMongoDB(t)
	defer clean()

	actual_col := db.GetCollectionName(o)

	if actual_col != expected_col {
		t.Fatalf("collection name is '%v', expected '%v'", actual_col, expected_col)
	}

	t.Logf("collection name is '%v'", actual_col)
}

func BenchmarkMongoDB_GetCollectionName(b *testing.B) {
	db, clean := newTestMongoDB(&testing.T{})
	defer clean()

	m := &mockRecord{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = db.GetCollectionName(m)
	}
}

func TestMongoDB_Find(t *testing.T) {

	var err error

	//setup
	db, clean := newTestMongoDB(t)
	defer clean()

	r1 := &mockRecord{}
	err = db.Save(r1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	r2 := &mockRecord{}
	r2.Object = r1.Object

	err = db.Find(r2)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	compareRecords(r1, r2, t)
}

func TestMongoDB_Remove(t *testing.T) {
	var err error

	//setup
	db, clean := newTestMongoDB(t)
	defer clean()

	r1 := &mockRecord{}
	err = db.Save(r1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	r2 := &mockRecord{}
	r2.Object = r1.Object

	//find the record
	err = db.Find(r2)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	//remove it
	err = db.Remove(r2)
	if err != nil {
		t.Fatalf("err =%v", err)
	}

	err = db.Find(r2)
	if err != ErrNotFound {
		t.Fatalf("error is '%v', expected '%v'", err, ErrNotFound)
	}
}

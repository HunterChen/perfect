package orm

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type MongoDB struct {
	Url         *url.URL
	Database    *mgo.Database
	name        string
	Session     *mgo.Session
	SessionInfo *mgo.BuildInfo
}

//NewDatabaseDriver
func NewMongoDBDriver(u *url.URL, name string) (Database, error) {
	return NewMongoDB(u, name), nil
}

func NewMongoDB(u *url.URL, name string) *MongoDB {
	if len(name) != 0 {
		u.Path = "/" + name
	} else {
		name = u.Path[1:]
	}

	return &MongoDB{
		Url:  u,
		name: name,
	}
}

func (db *MongoDB) Connect() (err error) {

	if db.Session == nil {
		//set connection properties
		dialinfo := &mgo.DialInfo{
			Addrs:    []string{db.Url.Host},
			FailFast: true,
			Database: db.name,
			Source:   db.name,
		}

		//set username/password info from the db url
		if db.Url.User != nil {
			dialinfo.Username = db.Url.User.Username()
			password, has_password := db.Url.User.Password()

			if has_password {
				dialinfo.Password = password
			}
		}

		//connect to MongoDB
		//db.Session, err = mgo.DialWithInfo(dialinfo)
		db.Session, err = mgo.DialWithInfo(dialinfo)
		if err != nil {
			return
		}
	}

	//create an mgo.Database object
	db.Database = db.Session.DB(db.Url.Path[1:])

	//fetch the session information
	info, err := db.Session.BuildInfo()
	if err != nil {
		return
	}

	db.SessionInfo = &info

	//set the session in safe mode
	db.Session.SetSafe(&mgo.Safe{})

	return
}

func (db *MongoDB) Disconnect() error {
	if db.Session != nil {
		db.Session.Close()
	}

	return nil
}

func (db *MongoDB) Name() string {
	return db.name
}

func (db *MongoDB) URL() *url.URL {
	return db.Url
}

func (db *MongoDB) SystemInformation() string {
	if db.SessionInfo != nil {
		return db.SessionInfo.SysInfo
	}

	return ""
}

func (db *MongoDB) Version() string {
	if db.SessionInfo != nil {
		return "MongoDB " + db.SessionInfo.Version
	}

	return ""
}

func (db *MongoDB) C(name string) Collection {
	if db.Database == nil {
		return nil
	}

	col := db.Database.C(name)

	if col == nil {
		return nil
	}

	return &MongoDBCollection{
		col,
	}
}

func (db *MongoDB) SetDebug(on bool) {
	mgo.SetDebug(on)
	if on {
		mgo.SetLogger(log.New(os.Stderr, "[db] ", log.LstdFlags))
	} else {
		mgo.SetLogger(nil)
	}
}

//returns a unique id (may be sequential)
func (db *MongoDB) UniqueId() string {
	return bson.NewObjectId().Hex()
}

func (db *MongoDB) GetCollectionName(r Record) string {
	name := strings.ToLower(reflect.ValueOf(r).Elem().Type().Name())

	//make sure collection names have plural form
	if !strings.HasSuffix(name, "s") {
		name += "s"
	}

	return name
}

func (db *MongoDB) Save(r Record) error {
	col_name := db.GetCollectionName(r)
	col := db.C(col_name)

	return col.Save(r)
}

func (db *MongoDB) Find(r Record) error {
	col_name := db.GetCollectionName(r)
	col := db.C(col_name)

	return col.Find(r)
}

func (db *MongoDB) Remove(r Record) error {
	col_name := db.GetCollectionName(r)
	col := db.C(col_name)

	return col.Remove(r)
}

func (db *MongoDB) DropCollection(r Record) error {
	col_name := db.GetCollectionName(r)
	col := db.C(col_name)

	return col.Drop()
}

func (db *MongoDB) Query(r Record) Query {
	col_name := db.GetCollectionName(r)
	col := db.C(col_name)

	return col.Query(r)
}

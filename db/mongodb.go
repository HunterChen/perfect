package db

import (
	"labix.org/v2/mgo"
	_ "labix.org/v2/mgo/bson"
	"net/url"
)

type MongoDB struct {
	Url      *url.URL
	Database *mgo.Database
	name     string
}

var (
	session     *mgo.Session
	sessionInfo *mgo.BuildInfo
)

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

	//create a single instance of the session object
	//which will get reused by subsequent calls to Connect()
	if session == nil {
		//set connection properties
		dialinfo := &mgo.DialInfo{
			Addrs:    []string{db.Url.Host},
			FailFast: true,
			Database: db.name,
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
		session, err = mgo.DialWithInfo(dialinfo)
		if err != nil {
			return
		}
	}

	//create an mgo.Database object
	db.Database = session.DB(db.Url.Path[1:])

	//fetch the session information
	info, err := session.BuildInfo()
	if err != nil {
		return
	}

	sessionInfo = &info

	//set the session in safe mode
	session.SetSafe(&mgo.Safe{})

	return
}

func (db *MongoDB) Disconnect() error {
	if session != nil {
		session.Close()
		session = nil
		sessionInfo = nil
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
	if sessionInfo != nil {
		return sessionInfo.SysInfo
	}

	return ""
}

func (db *MongoDB) Version() string {
	if sessionInfo != nil {
		return "MongoDB " + sessionInfo.Version
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

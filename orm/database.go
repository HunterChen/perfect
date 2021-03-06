package orm

import (
	"fmt"
	"io"
	"net/url"
)

//A database connection
type Database interface {
	Name() string
	URL() *url.URL
	SystemInformation() string
	Version() string

	Connect() error
	Disconnect() error

	C(string) Collection
	GetCollectionName(r Record) string
	DropCollection(Record) error

	SetDebug(bool)
	UniqueId() string

	Save(Record) error
	Find(Record) error
	Peek(Record) error
	Remove(Record) error

	Query(Record) Query

	NewLogger(col, prefix string) io.Writer
}

type Collection interface {
	Name() string
	Count() (n int, err error)
	Drop() error

	Save(Record) error
	Find(Record) error
	Peek(Record) error
	Remove(Record) error

	Query(interface{}) Query
}

type Query interface {
	Count() (int, error)
	One(Record) error
	Select(...string) Query
	Exclude(...string) Query
	All(interface{}) error
}

//All DB drivers must implement a NewDatabaseDriver function
//They accept a connection url and an optional database name.
type NewDatabaseDriver func(*url.URL, string) (Database, error)

//supported database adapters and URL schemes
var (
	schemes map[string]NewDatabaseDriver = map[string]NewDatabaseDriver{
		"mongodb": NewMongoDBDriver,
	}
)

//factory method to instantiate database drivers based on the url scheme
func NewDatabase(u *url.URL, name string) (Database, error) {
	factory, ok := schemes[u.Scheme]
	if !ok {
		err := fmt.Errorf("Unsupported database scheme '%s'", u.Scheme)
		return nil, err
	}

	return factory(u, name)
}

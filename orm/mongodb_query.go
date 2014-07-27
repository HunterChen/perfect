package orm

import (
	"labix.org/v2/mgo"
)

//several methods promoted from mgo.Query implement perfect/db.Query
type MongoDBQuery struct {
	*mgo.Query
}

func (q *MongoDBQuery) One(result Record) (err error) {
	err = q.Query.One(result)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}

	return
}

func (q *MongoDBQuery) All(result []Record) (err error) {
	err = q.Query.All(result)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}

	return
}

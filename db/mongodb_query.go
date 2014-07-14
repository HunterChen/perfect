package db

import (
	"labix.org/v2/mgo"
)

//several methods promoted from mgo.Query implement perfect/db.Query
type MongoDBQuery struct {
	*mgo.Query
}

func (q *MongoDBQuery) One(result Record) error {
	return q.Query.One(result)
}

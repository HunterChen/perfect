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

//TODO: write custom driver for mongodb so that 'result' could be of type []Record
func (q *MongoDBQuery) All(result interface{}) (err error) {
	err = q.Query.All(result)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}

	return
}

func (q *MongoDBQuery) Select(fields ...string) Query {

	var m = map[string]int{}
	for _, v := range fields {
		m[v] = 1
	}

	_ = q.Query.Select(m)

	return q
}

func (q *MongoDBQuery) Exclude(fields ...string) Query {
	var m = map[string]int{}
	for _, v := range fields {
		m[v] = 0
	}

	_ = q.Query.Select(m)

	return q
}

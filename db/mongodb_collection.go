package db

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var bogusDropError string = "ns not found"

var (
	ErrNotConnected error = errors.New("Not connected to database")
)

type MongoDBCollection struct {
	*mgo.Collection
}

func (col *MongoDBCollection) Name() string {
	if col.Collection == nil {
		return ""
	}

	return col.Collection.Name
}

func (col *MongoDBCollection) Drop() error {
	if col.Collection == nil {
		return nil
	}

	err := col.Collection.DropCollection()

	//ignore the "collection doesn't exist" error
	if err != nil && err.Error() == bogusDropError {
		return nil
	}

	return err
}

func (col *MongoDBCollection) Save(r Record) error {
	var err error

	if col.Collection == nil {
		return ErrNotConnected
	}

	id := r.GetDbId()

	if id == nil {
		id = bson.NewObjectId()
		r.SetDbId(id)
	}

	//update or insert a new object
	_, err = col.Collection.UpsertId(id, r)

	return err
}

func (col *MongoDBCollection) Find(r Record) error {
	if col.Collection == nil {
		return nil
	}

	return col.Collection.Find(r).One(r)
}

func (col *MongoDBCollection) Query(q interface{}) Query {
	if col.Collection == nil {
		return nil
	}

	return &MongoDBQuery{
		col.Collection.Find(q),
	}
}

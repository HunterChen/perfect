package orm

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
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
		err = col.Collection.Insert(r)
	} else {
		//update or insert a new object
		r.SetDbId(nil)
		_, err = col.Collection.UpsertId(id, bson.M{"$set": r})
		//don't return an error if the user attempted to save an empty object
		//the downside is that we can't insert an object if _id is set, we can
		//only update the document based on the _id
		if err != nil && strings.HasPrefix(err.Error(), "'$set' is empty") {
			err = nil
		}
		r.SetDbId(id)
	}

	return err
}

func (col *MongoDBCollection) Find(r Record) error {
	if col.Collection == nil {
		return nil
	}

	err := col.Collection.Find(r).One(r)

	if err == mgo.ErrNotFound {
		return ErrNotFound
	}

	return err
}

func (col *MongoDBCollection) Remove(r Record) error {
	if col.Collection == nil {
		return nil
	}

	err := col.Collection.Remove(r)
	if err == mgo.ErrNotFound {
		return ErrNotFound
	}

	return err
}

func (col *MongoDBCollection) Query(q interface{}) Query {
	if col.Collection == nil {
		return nil
	}

	return &MongoDBQuery{
		col.Collection.Find(q),
	}
}

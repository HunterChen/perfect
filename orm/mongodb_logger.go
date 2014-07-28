package orm

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type MongoDBLogger struct {
	Col    *mgo.Collection
	Prefix string
}

func (l *MongoDBLogger) Write(p []byte) (n int, err error) {
	if l.Col == nil {
		return 0, nil
	}

	n = len(p)

	doc := bson.M{
		"timestamp": time.Now(),
		"message":   l.Prefix + " " + string(p),
	}

	err = l.Col.Insert(doc)

	return
}

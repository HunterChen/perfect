package perfect

import (
	"github.com/vpetrov/perfect/db"
	_ "log"
)

const (
	SESSION_ID         = "SSESSIONID"
	SESSION_COLLECTION = "sessions"
)

//Represents a user's session. Id and _id are kept separate so that in
//the future, Id's can be regenerated on every request.
//Id and Authenticated are aliases for Values['id'] and Values['authenticated']
type Session struct {
	db.Object     `bson:",inline,omitempty" json:"-"`
	Id            *string            `bson:"id,omitempty" json:"id"`     //the publicly visible session id
	UserId        *string            `bson:"user_id,omitempty" json:"-"` //the user id this session is associated with
	Authenticated *bool              `bson:"authenticated" json:"-"`     //whether the user has logged in or not
	Values        *map[string]string `bson:"values" json:"-"`            //all other values go here
}

//creates a new Session object with no Id.
func NewSession(id string) *Session {
	return &Session{
		Id:            db.String(id),
		Authenticated: db.Bool(false),
		UserId:        nil,
		Values:        &map[string]string{},
	}
}

// Loads session info from the database
// returns nil if the session doesn't exist
func FindSession(id string, DB db.Database) (session *Session, err error) {
	session = &Session{
		Id: db.String(id),
	}

	err = DB.C(SESSION_COLLECTION).Find(session)

	//if the session doesn't exist, return error
	if err != nil {
		//use nil session to show that it was not found
		if err == ErrNotFound {
			err = nil
		}
		return nil, err
	}

	return
}

// Deletes itself from the database
func (s *Session) Delete(DB db.Database) (err error) {
	return DB.C(SESSION_COLLECTION).Remove(s)
}

// Stores the session in the database
func (s *Session) Save(DB db.Database) (err error) {
	return DB.C(SESSION_COLLECTION).Save(s)
}

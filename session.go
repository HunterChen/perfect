package perfect

import (
	"github.com/vpetrov/perfect/db"
	_ "log"
)

const (
	SESSION_ID         = "SSESSIONID"
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


package perfect

import (
	"github.com/vpetrov/perfect/orm"
	"net/http"
	"time"
)

const (
	SESSION_ID = "SSESSIONID"
)

var (
	SESSION_TIMEOUT time.Duration = time.Hour
)

//Represents a user's session. Id and _id are kept separate so that in
//the future, Id's can be regenerated on every request.
//Id and Authenticated are aliases for Values['id'] and Values['authenticated']
type Session struct {
	orm.Object    `bson:",inline,omitempty" json:"-"`
	Id            *string            `bson:"id,omitempty" json:"id"`           //the publicly visible session id
	ProfileId     *string            `bson:"profile_id,omitempty" json:"-"`    //the profile id this session is associated with
	Authenticated *bool              `bson:"authenticated,omitempty" json:"-"` //whether the user has logged in or not
	Values        *map[string]string `bson:"values,omitempty" json:"-"`        //all other values go here
}

//creates a new Session object with no Id.
func NewSession(id string) *Session {
	return &Session{
		Id:            orm.String(id),
		Authenticated: orm.Bool(false),
		ProfileId:     nil,
		Values:        &map[string]string{},
	}
}

func (session *Session) SetCookie(w http.ResponseWriter, r *Request) {
	//set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_ID,
		Value:    *session.Id,
		Path:     r.Module.MountPoint,
		Expires:  time.Now().Add(SESSION_TIMEOUT),
		Secure:   true,
		HttpOnly: true,
	})
}

func (session *Session) RemoveCookie(w http.ResponseWriter, r *Request) {
	//To delete the cookie, we set its value to some bogus string,
	//and the expiration to one second past the beginning of unix time.
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_ID,
		Value:    "Homer",
		Path:     r.Module.MountPoint,
		Expires:  time.Unix(1, 0),
		Secure:   true,
		HttpOnly: true,
	})
}

func (session *Session) ExtendCookie(w http.ResponseWriter, r *Request) {
	session.SetCookie(w, r)
}

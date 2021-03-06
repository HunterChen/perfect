package perfect

import (
	"github.com/vpetrov/perfect/json"
	"github.com/vpetrov/perfect/orm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//An interface for any type that can route Survana requests
type Router interface {
	Route(w http.ResponseWriter, r *Request)
}

//An HTTP request to a Survana component.
type RequestHandler func(http.ResponseWriter, *Request)

//a map of module-relative paths to their Handlers
type routeHandlers map[string]RequestHandler

//A struct that wraps http.Request and provides additional fields for all
//RequestHandlers to use. All methods from http.Request should be promoted
//for use with survana.Request.
type Request struct {
	*http.Request
	URL     *url.URL
	Module  *Module // the module that's handling the request
	session *Session
	profile *Profile
	Values  url.Values
}

// returns a new Request object
func NewRequest(r *http.Request, path string, module *Module) *Request {

	rurl := &url.URL{
		Scheme:   r.URL.Scheme,
		Opaque:   r.URL.Opaque,
		User:     r.URL.User,
		Host:     r.URL.Host,
		Path:     path,
		RawQuery: r.URL.RawQuery,
		Fragment: r.URL.Fragment,
	}

	req := &Request{
		Request: r,
		URL:     rurl,
		Module:  module,
	}

	if len(r.URL.RawQuery) != 0 {
		//parse query parameters, ignore any erorrs
		req.Values, _ = url.ParseQuery(r.URL.RawQuery)
	} else {
		req.Values = make(map[string][]string, 0)
	}

	return req
}

// returns either an existing session, or a new session
func (r *Request) Session() (*Session, error) {

	var (
		err error
	)

	db := r.Module.Db

	//if the session exists already, return it
	if r.session != nil {
		return r.session, nil
	}

	//get the session id cookie, if it exists
	session_id, _ := r.Cookie(SESSION_ID)

	session := &Session{Id: &session_id}

	//create a new session.
	err = db.Find(session)

	if err != nil {
		//if the session was not found, create a new one
		if err == orm.ErrNotFound {
			session = NewSession(MD5Sum(db.UniqueId()))
			err = db.Save(session)
			if err != nil {
				return nil, err
			}
		} else {
			//otherwise return the error
			return nil, err
		}
	}

	//cache the session object
	r.session = session

	return r.session, nil
}

func (r *Request) SetSession(s *Session) {
	r.session = s
}

//returns nil, nil if the profile was not found
func (r *Request) Profile() (*Profile, error) {
	var err error

	//if a profile already exists, return it
	if r.profile != nil {
		return r.profile, nil
	}

	//get the current session
	session, err := r.Session()
	if err != nil {
		return nil, err
	}

	//if there is no profile id, return 'not found'
	if session.ProfileId == nil {
		return nil, nil
	}

	//find the user profile by id (email)
	db := r.Module.Db
	profile := &Profile{Id: session.ProfileId}

	err = db.Find(profile)
	if err == orm.ErrNotFound {
		return nil, nil
	}

	//cache the profile object
	r.profile = profile

	return r.profile, nil
}

// returns the value of the cookie by name
func (r *Request) Cookie(name string) (value string, ok bool) {
	ok = false
	cookie, err := r.Request.Cookie(name)

	if cookie != nil && err == nil {
		ok = true
		value = cookie.Value
	}

	return
}

// Parses the request body as a JSON-encoded string
func (r *Request) ParseJSON(v interface{}) (err error) {
	return r.JSONBody(r.Request.Body, v)
}

func (r *Request) StringBody(body io.ReadCloser) (result string, err error) {
	bytes, err := r.BodyBytes(body)
	return string(bytes), err
}

func (r *Request) BodyBytes(body io.ReadCloser) (result []byte, err error) {
	// read the body
	result, err = ioutil.ReadAll(body)
	if err != nil {
		return
	}

	if len(result) == 0 {
		err = ErrEmptyRequest
		return
	}

	return
}

func (r *Request) JSONBody(body io.ReadCloser, v interface{}) (err error) {

	// read the body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}

	if len(data) == 0 {
		err = ErrEmptyRequest
		return
	}

	log.Println("API data:", string(data))

	//parse the JSON body
	err = json.Unmarshal(data, v)

	return
}

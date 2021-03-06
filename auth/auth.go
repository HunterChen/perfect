package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"log"
	"net/http"
)

//authentication strategy keys
const (
	BUILTIN = "built-in"
	OAUTH2  = "oauth2"
	LDAP    = "ldap"
	NIS     = "nis"
)

const (
	LOGIN_PATH = "/login"
)

var (
	auth_strategies map[string]NewStrategyFunc = map[string]NewStrategyFunc{
		BUILTIN: NewBuiltinStrategyFunc,
	}

	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
	ErrUsernameExists            = errors.New("Username already exists")
)

func New(config *Config) Strategy {

	newfunc, ok := auth_strategies[config.Type]
	if !ok {
		log.Fatalf("ERROR: Authentication type '%v' is not yet supported.", config.Type)
		return nil
	}

	if newfunc == nil {
		log.Fatalf("ERROR: No such authentication handler for '%v' authentication.", config.Type)
		return nil
	}

	return newfunc(config)
}

//TODO: use RWMutex in case this is called from different goroutines
func AddStrategy(name string, newfunc NewStrategyFunc) {
	auth_strategies[name] = newfunc
}

func Login(w http.ResponseWriter, r *perfect.Request) {
	var (
		profile_id *string
		err        error
	)

	//get the session
	session, err := r.Session()
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//if the user is already authenticated, redirect to home
	if *session.Authenticated {
		perfect.Redirect(w, r, "/")
		return
	}

	//TODO: make this work! (and delete the hard-coded built-in strategy!)
	//user, err := r.Module.Auth.Login(w, r);
	bauth := NewBuiltinStrategy(&Config{Type: BUILTIN})
	profile_id, err = bauth.Login(w, r)

	if err != nil {
		log.Println("login error:", err)
		perfect.JSONResult(w, r, false, err.Error())
		return
	}

	//mark the session as authenticated
	session.Authenticated = orm.Bool(true)

	//regenerate the session Id
	session.Id = orm.String(r.Module.Db.UniqueId())

	//set the current user profile id
	session.ProfileId = profile_id

	// update the session
	err = r.Module.Db.Save(session)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	session.SetCookie(w, r)

	//success
	perfect.JSONResult(w, r, true, r.Module.MountPoint+"/")
}

//Logs out a user.
//returns 204 No Content on success
//returns 500 Internal Server Error on failure
func logout(w http.ResponseWriter, r *perfect.Request) {
	session, err := r.Session()

	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	if !*session.Authenticated {
		perfect.NoContent(w)
		return
	}

	err = r.Module.Db.Remove(session)
	if err != nil && err != orm.ErrNotFound {
		perfect.Error(w, r, err)
		return
	}

	session.RemoveCookie(w, r)

	//return 204 No Content on success
	perfect.Redirect(w, r, "/")

	//note that the user has logged out
	go r.Module.Log.Printf("logout")
}

func generateRandomSalt(nbytes int) string {
	//allocate a slice of nbytes bytes
	random_bytes := make([]byte, nbytes)
	//read random data from a crypto rng
	rand.Read(random_bytes)

	return base64.StdEncoding.EncodeToString(random_bytes)
}

//salt prevents rainbow attacks. We insert some characters to block trivial attempts
//at guessing the salting function. That said, the attack could just look at the source
//code to find this function, but we still try to make it as hard as possible to guess it.
func salt(password, salt string) string {
	return "$" + salt + "::" + password + "::" + salt + "$"
}

func hash(password, password_salt string) []byte {
	//salt the password string
	salted_password := salt(password, password_salt)

	hash := sha512.New()

	return hash.Sum([]byte(salted_password))
}

//A handler that filters all requests that have not been authenticated
//returns 401 Unauthorized if the user's session hasn't been marked as authenticated
func Protect(handler perfect.RequestHandler) perfect.RequestHandler {
	return func(w http.ResponseWriter, r *perfect.Request) {
		//get the session
		session, err := r.Session()

		if err != nil {
			perfect.Error(w, r, err)
			return
		}

		//if the session hasn't been authorized, redirect
		if !*session.Authenticated {
			redirect_path := LOGIN_PATH

			w.Header().Set("X-Session-Expired", "1")

			//forward query params
			if len(r.URL.RawQuery) > 0 {
				redirect_path += "?" + r.URL.RawQuery
			}

			perfect.Redirect(w, r, redirect_path)
			return
		} else {
			//extend the amount of time the session is valid
			session.ExtendCookie(w, r)
		}

		//must be authenticated at this point
		handler(w, r)
	}
}

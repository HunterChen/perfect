package auth

import (
	"bytes"
	"errors"
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/orm"
	"log"
	"net/http"
)

const (
	SALT_ENTROPY             = 3
	BERR_INVALID_CREDENTIALS = "Invalid username or password"
)

type builtinUser struct {
	orm.Object `bson:",inline,omitempty"`
	Id         *string `bson:"id,omitempty"`
	Password   *[]byte `bson:"password,omitempty"`
	Salt       *string `bson:"salt,omitempty"`
	UserId     *string `bson:"user_id,omitempty"`
}

type BuiltinStrategy struct {
	Config *Config
}

func NewBuiltinStrategyFunc(config *Config) Strategy {
	return NewBuiltinStrategy(config)
}

func NewBuiltinStrategy(config *Config) *BuiltinStrategy {
	return &BuiltinStrategy{
		Config: config,
	}
}

func (b *BuiltinStrategy) Attach(module *perfect.Module) {
	app := module.Mux

	app.Get("/login", perfect.NotLoggedIn(b.LoginPage))
	app.Post("/login", perfect.NotLoggedIn(Login))

	//Registration is optional
	if b.Config.AllowRegistration {
		app.Get("/register", perfect.NotLoggedIn(b.RegistrationPage))
		app.Post("/register", perfect.NotLoggedIn(b.Register))
	}

	//setup the admin account
	if len(b.Config.Username) != 0 {
		user, profile, err := createBuiltinProfile(b.Config.Username, b.Config.Password, b.Config.Name, b.Config.Email, module.Db)
		if err != nil && err != ErrUsernameExists {
			log.Fatalf("Failed to create admin profile: %v", err)
		}

		log.Printf("Module administrator: username: %v email: %v\n", b.Config.Username, b.Config.Email)
	}
}

func (b *BuiltinStrategy) LoginPage(w http.ResponseWriter, r *perfect.Request) {
	r.Module.RenderTemplate(w, r, "auth/builtin/login", b.Config)
}

func (b *BuiltinStrategy) RegistrationPage(w http.ResponseWriter, r *perfect.Request) {
	r.Module.RenderTemplate(w, r, "auth/builtin/register", nil)
}

func (b *BuiltinStrategy) Login(w http.ResponseWriter, r *perfect.Request) (profile_id *string, err error) {

	//this is why each strategy needs to be able to render its
	//login screens, so that it can ask for custom fields.
	//here we have a simple username/password combo, but the
	//other strategies could show various options based on the
	//auth configuration
	data := make(map[string]string)

	err = r.ParseJSON(&data)
	if err != nil {
		log.Println(err)
		return
	}

	username, ok1 := data["username"]
	password, ok2 := data["password"]

	if !ok1 || !ok2 || len(username) == 0 || len(password) == 0 {
		err = errors.New("Invalid request")
		return
	}

	user := &builtinUser{Id: &username}

	//find this user in the built-in user database
	err = r.Module.Db.Find(user)
	if err == orm.ErrNotFound {
		log.Printf("No such builtin user: %v", username)
		return nil, ErrInvalidUsernameOrPassword
	} else if err != nil {
		return nil, err
	}

	sha512_password := hash(password, *user.Salt)

	//wrong password?
	if !bytes.Equal(sha512_password, *user.Password) {
		err = errors.New(BERR_INVALID_CREDENTIALS)
		return
	}

	return user.UserId, nil
}

func (b *BuiltinStrategy) Register(w http.ResponseWriter, r *perfect.Request) {

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

	data := make(map[string]string)

	err = r.ParseJSON(&data)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	username, ok1 := data["username"]
	password, ok2 := data["password"]
	name, ok3 := data["name"]
	email, ok4 := data["email"]

	//TODO: this needs to be refactored into something better
	if !ok1 || !ok2 || !ok3 || !ok4 || len(username) == 0 || len(password) == 0 || len(name) == 0 || len(email) == 0 {
		perfect.JSONResult(w, r, false, "Please complete all fields")
		return
	}

	_, _, err = createBuiltinProfile(username, password, email, name, r.Module.Db)
	if err != nil {
		perfect.JSONResult(w, r, true, r.Module.MountPoint+"/")
		return
	}

	perfect.JSONResult(w, r, false, err)
}

//default logout
func (b *BuiltinStrategy) Logout(w http.ResponseWriter, r *perfect.Request) {
	//call the default logout func in auth.go
	logout(w, r)
}

func createBuiltinProfile(username, password, name, email string, db orm.Database) (user *builtinUser, profile *perfect.Profile, err error) {
	user = &builtinUser{Id: &username}

	//query the database to check if the username exists
	err = db.Find(user)
	if err != orm.ErrNotFound {
		//make sure users can't register duplicate usernames
		if err == nil {
			err = ErrUsernameExists
		}
		return
	}

	//create a perfect user (profile)
	profile = perfect.NewProfile(email, name)
	profile.AuthType = orm.String(BUILTIN)
	err = db.Save(profile)
	if err != nil {
		log.Println(err)
		return
	}

	//hash the password
	password_salt := generateRandomSalt(SALT_ENTROPY)
	password_hash := hash(password, password_salt)

	//create an entry to store auth details
	user.Password = &password_hash
	user.Salt = &password_salt
	user.UserId = profile.Id //TODO: change to ProfileId
	err = db.Save(user)
	if err != nil {
		log.Println(err)
		return
	}

	return user, profile, nil
}

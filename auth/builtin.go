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
	orm.Object `bson:",inline,omitempty" json:"-"`
	Id         *string `bson:"id,omitempty" json:"id,omitempty"`
	Password   *[]byte `bson:"password,omitempty" json:"-"`
	Salt       *string `bson:"salt,omitempty" json:"-"`
	ProfileId  *string `bson:"profile_id,omitempty" json:"profile_id,omitempty"`
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
	module.Get("/login", perfect.NotLoggedIn(b.LoginPage))
	module.Post("/login", perfect.NotLoggedIn(Login))

	//Registration is optional
	if b.Config.AllowRegistration {
		module.Get("/register", perfect.NotLoggedIn(b.RegistrationPage))
		module.Post("/register", perfect.NotLoggedIn(b.Register))
	}

	if len(b.Config.Username) != 0 {
		user, profile, err := b.setupAdminAccount(module.Db)
		if err != nil {

		}
		log.Printf("Module administrator: username: %v email: %v\n", *user.Id, *profile.Id)
	} else {
		log.Printf("WARNING: No authentication details found for module '%v'", module.Name)
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

	return user.ProfileId, nil
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
	user.ProfileId = profile.Id
	err = db.Save(user)
	if err != nil {
		log.Println(err)
		return
	}

	return user, profile, nil
}

func (b *BuiltinStrategy) setupAdminAccount(db orm.Database) (user *builtinUser, profile *perfect.Profile, err error) {

	//setup the admin account
	user = &builtinUser{Id: orm.String(b.Config.Username)}
	err = db.Find(user)
	if err != nil && err != orm.ErrNotFound {
		return
	}

	profile = &perfect.Profile{}
	//search for this profile
	if user.ProfileId != nil {
		profile.Id = user.ProfileId
		err = db.Find(profile)
		if err != nil && err != orm.ErrNotFound {
			return
		}
	}

	//update the profile
	profile.Id = orm.String(b.Config.Email)
	profile.Name = orm.String(b.Config.Name)
	profile.AuthType = orm.String(b.Config.Type)
	err = db.Save(profile)
	if err != nil {
		return
	}

	//hash the password
	password_salt := generateRandomSalt(SALT_ENTROPY)
	password_hash := hash(b.Config.Password, password_salt)
	//update username
	user.Password = &password_hash
	user.Salt = &password_salt
	user.ProfileId = profile.Id
	err = db.Save(user)
	if err != nil {
		return
	}

	return user, profile, nil
}

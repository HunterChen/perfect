package auth

import (
        "github.com/vpetrov/perfect"
        "net/http"
       )

type Strategy interface {
    Attach(module *perfect.Module)

    LoginPage(w http.ResponseWriter, r *perfect.Request)
    RegistrationPage(w http.ResponseWriter, r *perfect.Request)

    Login(w http.ResponseWriter, r *perfect.Request) (profile_id string, err error)
    Register(w http.ResponseWriter, r *perfect.Request)
    Logout(w http.ResponseWriter, r *perfect.Request)
}

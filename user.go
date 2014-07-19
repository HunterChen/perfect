package perfect

import (
	"github.com/vpetrov/perfect/db"
	_ "log"
)

type User struct {
	db.Object `bson:",inline,omitempty" json:"-"`
	Id        *string   `bson:"id,omitempty" json:"id,omitempty"`
	Name      *string   `bson:"name,omitempty" json:"name,omitempty"`
	Groups    *[]string `bson:"groups,omitempty" json:"groups,omitempty"`
	AuthType  *string   `bson:"auth_type,omitempty" json:"auth_type,omitempty"`
}

func NewUser(email, name string) *User {
	return &User{
		Id:   db.String(email),
		Name: db.String(name),
	}
}


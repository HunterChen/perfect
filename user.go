package perfect

import (
	"github.com/vpetrov/perfect/orm"
)

type User struct {
	orm.Object `bson:",inline,omitempty" json:"-"`
	Id         *string   `bson:"id,omitempty" json:"id,omitempty"`
	Name       *string   `bson:"name,omitempty" json:"name,omitempty"`
	Groups     *[]string `bson:"groups,omitempty" json:"groups,omitempty"`
	AuthType   *string   `bson:"auth_type,omitempty" json:"auth_type,omitempty"`
}

func NewUser(email, name string) *User {
	return &User{
		Id:   orm.String(email),
		Name: orm.String(name),
	}
}

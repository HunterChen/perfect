package orm

//An object that can be stored in a database; implements the Record interface
type Object struct {
	Id interface{} `bson:"_id,omitempty" json:"-"`
}

func (o *Object) GetDbId() interface{} {
	return o.Id
}

func (o *Object) SetDbId(id interface{}) {
	o.Id = id
}

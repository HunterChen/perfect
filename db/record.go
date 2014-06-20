package db

type Record interface {
	GetDbId() interface{}
	SetDbId(interface{})
}

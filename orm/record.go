package orm

type Record interface {
	GetDbId() interface{}
	SetDbId(interface{})
}

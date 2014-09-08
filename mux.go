package perfect

type Mux interface {
	Router
	Handle(method, path string, handler RequestHandler)
	Get(path string, handler RequestHandler)
	Post(path string, handler RequestHandler)
	Put(path string, handler RequestHandler)
	Delete(path string, handler RequestHandler)
	Head(path string, handler RequestHandler)

	Static(path string)
}

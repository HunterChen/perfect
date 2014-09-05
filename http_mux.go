package perfect

import (
	"log"
	"net/http"
	"strings"
)

//a map of module-relative paths to their Handlers
type RouteHandlers map[string]RequestHandler

// This Request Mux does not use a mutex because we're not anticipating the need
// to change the routes at run time, as requests are served.
// Even if modules are re-mounted, they should be first instantiated,
// then enabled. Should this change, an RWMutex will be necessary.
type HTTPMux struct {
	Handlers       map[string]RouteHandlers
	StaticPrefix   string
	HasStaticFiles bool
}

//returns a new Mux
func NewHTTPMu() *Mux {
	return &Mux{
		Handlers:       make(map[string]RouteHandlers, 0),
		StaticPrefix:   "",
		HasStaticFiles: false,
	}
}

// finds and invokes the Handlers for the given request
func (h *Mux) Route(w http.ResponseWriter, r *Request) {
	handler := h.FindHandler(r)

	if handler == nil {
		http.NotFound(w, r.Request)
		return
	}

	//TODO: remove these 2 lines when optimizing performance
	name, file, line := HandlerInfo(handler)
	log.Printf("%s %s%s -> %s at %s:%d\n", r.Method, r.Module.MountPoint, r.URL, name, file, line)

	//invoke the handler
	handler(w, r)
}

//checks whether a request is for a static resource
func (h *Mux) isStatic(r *Request) bool {
	return h.HasStaticFiles && strings.HasPrefix(r.URL.Path, h.StaticPrefix)
}

//generic method that registers a handler for a path and http method
func (h *Mux) Handle(method string, path string, handler RequestHandler) {

	Handlers, ok := h.Handlers[method]

	if !ok {
		h.Handlers[method] = make(RouteHandlers, 0)
		Handlers = h.Handlers[method]
	}

	Handlers[path] = handler

	log.Println("[mux]", method, path, handler)
}

//sets the static path
func (h *Mux) Static(path string) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}

	h.StaticPrefix = path
	h.HasStaticFiles = true
}

//a request handler for static resources
func (h *Mux) StaticHandler(w http.ResponseWriter, r *Request) {
	http.ServeFile(w, r.Request, r.Module.Path+r.URL.Path)
}

//registers a GET request handler
func (h *Mux) Get(path string, handler RequestHandler) {
	h.Handle("GET", path, handler)
}

//registers a POST request handler
func (h *Mux) Post(path string, handler RequestHandler) {
	h.Handle("POST", path, handler)
}

//registers a PUT request handler
func (h *Mux) Put(path string, handler RequestHandler) {
	h.Handle("PUT", path, handler)
}

//registers a DELETE request handler
func (h *Mux) Delete(path string, handler RequestHandler) {
	h.Handle("DELETE", path, handler)
}

//registers a HEAD request handler
func (h *Mux) Head(path string, handler RequestHandler) {
	h.Handle("HEAD", path, handler)
}

//returns a static/dynamic request handler for the given request
func (h *Mux) FindHandler(r *Request) RequestHandler {

	//if the request is for a static resource, return the
	//static request handler
	if h.HasStaticFiles && h.isStatic(r) {
		return h.StaticHandler
	}

	//GET, POST, PUT
	route_handler, ok := h.Handlers[r.Request.Method]

	if !ok {
		return nil
	}

	//Find the handler
	handler, ok := route_handler[r.URL.Path]

	if !ok {
		return nil
	}

	return handler
}

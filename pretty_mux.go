package perfect

import (
	"log"
)

type PrettyMux struct {
	*Mux
}

func NewPrettyMux() *PrettyMux {
	return &PrettyMux{
		Mux: NewMux(),
	}
}

//generic method that registers a handler for a path and http method
func (h *PrettyMux) Handle(method string, expr string, handler RequestHandler) {

	Handlers, ok := h.Handlers[method]

	if !ok {
		h.Handlers[method] = make(RouteHandlers, 0)
		Handlers = h.Handlers[method]
	}

	Handlers[expr] = handler

	log.Println("[mux]", method, expr, handler)
}

//registers a GET request handler
func (h *PrettyMux) Get(expr string, handler RequestHandler) {
	h.Handle("GET", expr, handler)
}

//registers a POST request handler
func (h *PrettyMux) Post(path string, handler RequestHandler) {
	h.Handle("POST", path, handler)
}

//registers a PUT request handler
func (h *PrettyMux) Put(path string, handler RequestHandler) {
	h.Handle("PUT", path, handler)
}

//registers a DELETE request handler
func (h *PrettyMux) Delete(path string, handler RequestHandler) {
	h.Handle("DELETE", path, handler)
}

//registers a HEAD request handler
func (h *PrettyMux) Head(path string, handler RequestHandler) {
	h.Handle("HEAD", path, handler)
}

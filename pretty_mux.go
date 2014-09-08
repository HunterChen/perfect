package perfect

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type prettyRoute struct {
	path     string
	elements []string
	handler  RequestHandler
}

func (r *prettyRoute) Match(path_elements []string) (bool, *url.Values) {
	//quick check to see if the number of elements are the same
	if len(r.elements) != len(path_elements) {
		return false, nil
	}

	values := &url.Values{}

	for i, e := range r.elements {
		//check named parameters
		if strings.HasPrefix(e, ":") {
			parts := strings.SplitN(e, ":", 2)
			//only store the value of the named parameter if it exists.
			//this also means that /:/ is valid syntax and does not result
			//in a named parameter
			if len(parts) > 1 && len(parts[1]) > 0 {
				values.Add(parts[1], path_elements[i])
			}
			continue
		}

		if e != path_elements[i] {
			return false, nil
		}
	}

	return true, values
}

type PrettyMux struct {
	StaticPrefix   string
	hasStaticFiles bool

	//map [HTTP_METHOD] list of pretty routes
	routes map[string][]prettyRoute
}

func NewPrettyMux() *PrettyMux {
	return &PrettyMux{
		routes: make(map[string][]prettyRoute),
	}
}

//generic method that registers a handler for a path and http method
func (pm *PrettyMux) Handle(method string, expr string, handler RequestHandler) {

	_, ok := pm.routes[method]

	if !ok {
		pm.routes[method] = make([]prettyRoute, 0)
	}

	expr = strings.TrimSuffix(expr, "/")

	route := prettyRoute{
		path:     expr,
		elements: strings.Split(expr, "/"),
		handler:  handler,
	}

	pm.routes[method] = append(pm.routes[method], route)

	log.Println("[pretty mux]", method, expr, handler)
}

//registers a GET request handler
func (pm *PrettyMux) Get(expr string, handler RequestHandler) {
	pm.Handle("GET", expr, handler)
}

//registers a POST request handler
func (pm *PrettyMux) Post(path string, handler RequestHandler) {
	pm.Handle("POST", path, handler)
}

//registers a PUT request handler
func (pm *PrettyMux) Put(path string, handler RequestHandler) {
	pm.Handle("PUT", path, handler)
}

//registers a DELETE request handler
func (pm *PrettyMux) Delete(path string, handler RequestHandler) {
	pm.Handle("DELETE", path, handler)
}

//registers a HEAD request handler
func (pm *PrettyMux) Head(path string, handler RequestHandler) {
	pm.Handle("HEAD", path, handler)
}

//returns a static/dynamic request handler for the given request
func (pm *PrettyMux) FindHandler(r *Request) RequestHandler {

	//if the request is for a static resource, return the
	//static request handler
	if pm.hasStaticFiles && pm.isStatic(r) {
		return pm.StaticHandler
	}

	//GET, POST, PUT
	routes, ok := pm.routes[r.Request.Method]

	if !ok {
		return nil
	}

	path_elements := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	for _, route := range routes {
		matches, values := route.Match(path_elements)
		if matches {
			//prepend all values to request
			//prepending mimics the order of these values in the URL, in case
			//of repeating parameters (i.e. /:id/:id -> /1/2?id=3) -> {"id": [1,2,3]}

			for p, v := range *values {
				r.Values[p] = append(v, r.Values[p]...)
			}

			return route.handler
		}
	}

	return nil
}

//checks whether a request is for a static resource
func (pm *PrettyMux) isStatic(r *Request) bool {
	return pm.hasStaticFiles && strings.HasPrefix(r.URL.Path, pm.StaticPrefix)
}

//a request handler for static resources
func (pm *PrettyMux) StaticHandler(w http.ResponseWriter, r *Request) {
	http.ServeFile(w, r.Request, r.Module.Path+r.URL.Path)
}

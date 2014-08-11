package perfect

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

//TODO: remove this or make it configurable
const GITHUB_URL = "https://api.github.com/repos/vpetrov/survana/issues?access_token=a8202d411154a49dab3a42b6e46fc51549ce651a"

//decides which Module handles which Request
type ModuleMux struct {
	lock    sync.RWMutex
	modules map[string]*Module
}

//Returns a new ModuleMux
func NewModuleMux() *ModuleMux {
	return &ModuleMux{
		modules: make(map[string]*Module, 0),
	}
}

//Finds the requested module and hands off the request to its router
func (mux *ModuleMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	//detect which module this request should go to
	mount_point, rurl := mux.GetModule(r.URL.Path)

	//lock modules mutex for reading (ensures that the map won't be changed
	//while we're reading from it)
	mux.lock.RLock()
	//fetch the module
	module, ok := mux.modules[mount_point]
	mux.lock.RUnlock()

	if !ok {
		//fetch the "/" module
		mux.lock.RLock()
		module, ok = mux.modules["/"]
		mux.lock.RUnlock()
		// if no default module found, return a 500 Internal Server Error
		if !ok {
			http.Error(w,
				"Internal Server Error",
				http.StatusInternalServerError)
			return
		}
	}

	//create an application-specific request object
	request := NewRequest(r, rurl, module)

	//report panics as issues on Github (with debug info)
	defer func() {
		if r_err := recover(); r_err != nil {

			var error_string string

			switch c := r_err.(type) {
			case string:
				error_string = c
			case error:
				error_string = c.Error()
			default:
				error_string = "unknown error object"
			}

			//convert the http request to a json string
			json_request, err := json.MarshalIndent(request, "  ", "")
			if err != nil {
				log.Println("DEBUG: Failed to encode request as JSON:", err)
				panic(r_err)
			}
			//create the report: error, stack, request
			report := &GithubIssue{
				Title:  "Automatic Panic Report (" + module.Name + ")",
				Body:   error_string + "\n\n Stack:\n" + string(debug.Stack()) + "\n\n Request:\n" + string(json_request),
				Labels: []string{"auto", "panic", module.Name},
			}

			//convert the report to JSON
			json_report, err := json.Marshal(report)
			buf := bytes.NewBuffer(json_report)
			github_response, err := http.Post(GITHUB_URL, "application/json", buf)
			if err != nil {
				log.Println("DEBUG: Failed to create new GitHub issue:", err)
				panic(r_err)
			}

			json_response, err := ioutil.ReadAll(github_response.Body)
			if err != nil {
				log.Println("DEBUG: Failed to read GitHub's response:", err)
				panic(r_err)
			}

			response := &GithubResponse{}
			//read Github's reply
			err = json.Unmarshal(json_response, response)
			if err != nil {
				log.Println("DEBUG: Failed to read reply from Github:", err)
				panic(r_err)
			}

			//set issue info as special headers
			if response.IssueNumber > 0 {
				w.Header().Set("X-Survana-IssueNumber", string(response.IssueNumber))
				w.Header().Set("X-Survana-IssueUrl", response.HtmlUrl)
			}

			//send the original error to the client
			Error(w, request, r_err.(error))
		}
	}()

	//route the request
	module.Route(w, request)

	if module.Log != nil {
		go module.Log.Printf("%s|%s|%s|%s", r.Method, r.URL.Path, time.Since(startTime).String(), r.UserAgent())
	}
}

//Registers a new Module for a URL path
func (mux *ModuleMux) Mount(m *Module, path string) {
	//atempt to lock the modules mutex before we write to the map
	mux.lock.Lock()
	// add the module pointer to the map
	mux.modules[path] = m
	// unlock the mutex
	mux.lock.Unlock()

	log.Println("Mounting ", m.Name, "on", path)

	m.MountPoint = path
}

//Unregisters the module that handles the path
func (mux *ModuleMux) Unmount(path string) {
	mux.lock.Lock()
	//remove the module from the path
	delete(mux.modules, path)
	mux.lock.Unlock()
}

// searches for a module by the path it's been mounted on
func (mux *ModuleMux) GetModule(path string) (module, mpath string) {

	if len(path) <= 1 {
		return "/", path
	}

	// find the second "/" in path
	slash := strings.Index(path[1:], "/") + 1

	//if no slashes were found, return original path and default module
	if slash == 0 {
		return "/", path
	}

	// the module is the first item between the 2 slashes
	module = path[:slash]

	// the module's URL 'path' is everything that follows it
	mpath = path[slash:]

	return
}

//the default module mux
var Modules *ModuleMux = NewModuleMux()

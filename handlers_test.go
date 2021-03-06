package perfect

import (
	"encoding/json"
	"errors"
	"github.com/vpetrov/perfect/orm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNotLoggedIn(t *testing.T) {
	module_mount_point := "/testing"
	request_method := "GET"
	request_path := "/test"
	query_string := "?arg1=val1"
	session_id := "123ABC"

	session := NewSession(session_id)

	module := &Module{
		MountPoint: module_mount_point,
	}

	request_url, err := url.Parse("http://localhost" + module.MountPoint + request_path + query_string)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	response := httptest.NewRecorder()

	http_request := &http.Request{
		Method: request_method,
		URL:    request_url,
		Header: http.Header{},
	}

	request := NewRequest(http_request, request_path, module)

	request.SetSession(session)

	handler := func(w http.ResponseWriter, r *Request) {
		if *session.Authenticated {
			t.Errorf("session.Authenticated is %v, did not expect the handler to be called", session.Authenticated)
		}

		w.WriteHeader(http.StatusOK)
	}

	auth_handler := NotLoggedIn(handler)

	//call the auth handler
	auth_handler(response, request)

	if response.Code == http.StatusSeeOther {
		t.Errorf("response.Code is %v, expected %v (http.StatusOK)", response.Code, http.StatusSeeOther)
	}

	//perform another test, but this time the session is going to be authenticated
	session.Authenticated = orm.Bool(true)
	request.SetSession(session)
	response = httptest.NewRecorder()

	handler2_called := false

	handler2 := func(w http.ResponseWriter, r *Request) {
		handler2_called = true
	}

	auth_handler = NotLoggedIn(handler2)

	auth_handler(response, request)

	if handler2_called {
		t.Errorf("expected handler2 to not be called")
	}

	if response.Code != http.StatusSeeOther {
		t.Errorf("response.Code is %v, expected %v (http.StatusSeeOther)", response.Code, http.StatusSeeOther)
	}
}

func TestNotFound(t *testing.T) {
	response := httptest.NewRecorder()

	NotFound(response)

	if response.Code != http.StatusNotFound {
		t.Errorf("responses.Status is %v, expected %v (http.StatusNotFound)", response.Code, http.StatusNotFound)
	}

	if response.Body.Len() == 0 {
		t.Errorf("response.Body.Len() is %v, expected non-zero", response.Body.Len())
	}
}

func TestNoContent(t *testing.T) {
	response := httptest.NewRecorder()

	NoContent(response)

	if response.Code != http.StatusNoContent {
		t.Errorf("responses.Status is %v, expected %v (http.StatusNoContent)", response.Code, http.StatusNoContent)
	}

	if response.Body.Len() != 0 {
		t.Errorf("response.Body.Len() is %v, expected zero", response.Body.Len())
	}
}

func TestError(t *testing.T) {
	request := &Request{
		Module:  &Module{},
		Request: &http.Request{},
	}
	response := httptest.NewRecorder()
	err := errors.New("test")

	Error(response, request, err)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("response.Code is %v, expected %v (http.StatusInternalServerError)", response.Code, http.StatusNoContent)
	}

	if response.Body.Len() == 0 {
		t.Errorf("response.Body.Len() is %v, expected non-zero", response.Body.Len())
	}
}

func TestRedirectSimple(t *testing.T) {

	module_mount_point := "/testing"
	request_method := "GET"
	request_path := "/test"
	redirect_path := "/test2"
	query_string := "?arg1=val1"

	module := &Module{
		MountPoint: module_mount_point,
	}

	request_url, err := url.Parse("http://localhost" + module.MountPoint + request_path + query_string)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	response := httptest.NewRecorder()

	http_request := &http.Request{
		Method: request_method,
		URL:    request_url,
		Header: http.Header{},
	}

	request := NewRequest(http_request, request_path, module)

	Redirect(response, request, redirect_path)

	//response.Code has to be StatusSeeOther
	if response.Code != http.StatusSeeOther {
		t.Errorf("response.Code is %v, expected %v (http.StatusSeeOther)", response.Code, http.StatusSeeOther)
	}

	location := response.HeaderMap.Get("Location")

	//the redirect must be relative to the module mount point
	if !strings.HasPrefix(location, module.MountPoint) {
		t.Errorf("Redirect Location header value is %v, expected this path to be prefixed with the module's mount point, %v", location, module.MountPoint)
	}

	//the redirect location must contain the redirect path
	if strings.Index(location, redirect_path) < len(module.MountPoint) {
		t.Errorf("Redirect Location header value is %v, expected this path to contain the redirect path %v", location, redirect_path)
	}

	//verify that X-Survana-Redirect isn't set
	_, ok := response.HeaderMap["X-Survana-Redirect"]

	if ok {
		t.Errorf("X-Survana-Redirect header is '%v', expected this header value to not be set.", location)
	}
}

func TestRedirectWithXHR(t *testing.T) {
	module_mount_point := "/testing"
	request_method := "GET"
	request_path := "/test"
	redirect_path := "/test2"
	query_string := "?arg1=val1"

	module := &Module{
		MountPoint: module_mount_point,
	}

	request_url, err := url.Parse("http://localhost" + module.MountPoint + request_path + query_string)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	response := httptest.NewRecorder()

	http_request := &http.Request{
		Method: request_method,
		URL:    request_url,
		Header: http.Header{},
	}

	request := NewRequest(http_request, request_path, module)

	//set X-Requested-With to XMLHttpRequest
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	Redirect(response, request, redirect_path)

	//verify that X-Survana-Redirect is set
	x_survana_redirect_header, ok := response.HeaderMap["X-Survana-Redirect"]

	if !ok {
		t.Errorf("X-Survana-Redirect header not present, expected it to be set")
	}

	if len(x_survana_redirect_header) == 0 {
		t.Errorf("X-Survana-Redirect header is empty, expected it to have at least 1 value")
	}

	x_survana_redirect := x_survana_redirect_header[0]

	//the redirect must be relative to the module mount point
	if !strings.HasPrefix(x_survana_redirect, module.MountPoint) {
		t.Errorf("X-Survana-Redirect header value is %v, expected this path to be prefixed with the module's mount point, %v", x_survana_redirect, module.MountPoint)
	}

	//the redirect location must contain the redirect path
	if strings.Index(x_survana_redirect, redirect_path) < len(module.MountPoint) {
		t.Errorf("X-Survana-Redirect header value is %v, expected this path to contain the redirect path %v", x_survana_redirect, redirect_path)
	}
}

func TestXHRRedirect(t *testing.T) {

	module_mount_point := "/testing"
	request_method := "GET"
	request_path := "/test"
	redirect_path := module_mount_point + "/test2"
	query_string := "?arg1=val1"

	module := &Module{
		MountPoint: module_mount_point,
	}

	request_url, err := url.Parse("http://localhost" + module.MountPoint + request_path + query_string)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	response := httptest.NewRecorder()

	http_request := &http.Request{
		Method: request_method,
		URL:    request_url,
		Header: http.Header{},
	}

	request := NewRequest(http_request, request_path, module)

	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	XHRRedirect(response, request, redirect_path)

	//verify that X-Survana-Redirect is set
	x_survana_redirect_header, ok := response.HeaderMap["X-Survana-Redirect"]

	if !ok {
		t.Errorf("X-Survana-Redirect header not present, expected it to be set")
	}

	if len(x_survana_redirect_header) == 0 {
		t.Errorf("X-Survana-Redirect header is empty, expected it to have at least 1 value")
	}

	x_survana_redirect := x_survana_redirect_header[0]

	//make sure the X-Survana-Redirect value is the same as the redirect_path
	if x_survana_redirect != redirect_path {
		t.Errorf("X-Survana-Redirect header value is %v, expected %v", x_survana_redirect, redirect_path)
	}
}

func TestJSONResult(t *testing.T) {
	request := &Request{
		Module:  &Module{},
		Request: &http.Request{},
	}
	response := httptest.NewRecorder()
	expected_response := &JSONResponse{
		Success: true,
		Message: "This is a test",
	}

	JSONResult(response, request, expected_response.Success, expected_response.Message)

	if response.Body.Len() == 0 {
		t.Errorf("response.Body.Len() is %v, expected non-zero", response.Body.Len())
	}

	//unmarshal the response
	actual_response := JSONResponse{}
	err := json.Unmarshal(response.Body.Bytes(), &actual_response)
	if err != nil {
		t.Errorf("err = %v", err)
	}

	if actual_response.Success != expected_response.Success {
		t.Errorf("response success value is %v, expected %v", actual_response.Success, expected_response.Success)
	}

	if actual_response.Message != expected_response.Message {
		t.Errorf("response message is %v, expected %v", actual_response.Message, expected_response.Message)
	}
}

func TestBadRequest(t *testing.T) {
	response := httptest.NewRecorder()

	BadRequest(response)

	if response.Code != http.StatusBadRequest {
		t.Errorf("response.Code is %v, expected %v (http.StatusBadRequest)", response.Code, http.StatusBadRequest)
	}

	if response.Body.Len() == 0 {
		t.Errorf("response.Body.Len() is %v, expected non-zero", response.Body.Len())
	}
}

func TestUnauthorized(t *testing.T) {
	response := httptest.NewRecorder()
	err := errors.New("Access is denied")

	Unauthorized(response, err)

	if response.Code != http.StatusUnauthorized {
		t.Errorf("response.Code is %v, expected %v (http.StatusUnauthorized)", response.Code, http.StatusUnauthorized)
	}

	if response.Body.Len() == 0 {
		t.Errorf("response.Body.Len() is %v, expected non-zero", response.Body.Len())
	}
}

func TestHandlerInfo(t *testing.T) {
	handler := func(_ http.ResponseWriter, _ *Request) {
		//dummy
	}

	name, file, line := HandlerInfo(handler)

	if len(name) == 0 {
		t.Errorf("handler name is %v, expected a non-empty string", name)
	}

	if len(file) == 0 {
		t.Errorf("len(file) is %v, expected non-zero", len(file))
	}

	if line <= 0 {
		t.Errorf("line is %v, expected a number greater than zero", line)
	}
}

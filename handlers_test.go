package perfect

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
    DB "github.com/vpetrov/perfect/db"
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

	response := NewMockResponse()

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

	if response.Status == http.StatusSeeOther {
		t.Errorf("response.Status is %v, expected %v (http.StatusOK)", response.Status, http.StatusSeeOther)
	}

	//perform another test, but this time the session is going to be authenticated
	session.Authenticated = DB.Bool(true)
	request.SetSession(session)
	response = NewMockResponse()

	handler2_called := false

	handler2 := func(w http.ResponseWriter, r *Request) {
		handler2_called = true
	}

	auth_handler = NotLoggedIn(handler2)

	auth_handler(response, request)

	if handler2_called {
		t.Errorf("expected handler2 to not be called")
	}

	if response.Status != http.StatusSeeOther {
		t.Errorf("response.Status is %v, expected %v (http.StatusSeeOther)", response.Status, http.StatusSeeOther)
	}
}

func TestNotFound(t *testing.T) {
	response := NewMockResponse()

	NotFound(response)

	if response.Status != http.StatusNotFound {
		t.Errorf("responses.Status is %v, expected %v (http.StatusNotFound)", response.Status, http.StatusNotFound)
	}

	if len(response.Data) == 0 {
		t.Errorf("len(response.Data) is %v, expected non-zero", len(response.Data))
	}
}

func TestNoContent(t *testing.T) {
	response := NewMockResponse()

	NoContent(response)

	if response.Status != http.StatusNoContent {
		t.Errorf("responses.Status is %v, expected %v (http.StatusNoContent)", response.Status, http.StatusNoContent)
	}

	if len(response.Data) != 0 {
		t.Errorf("len(response.Data) is %v, expected zero", len(response.Data))
	}
}

func TestError(t *testing.T) {
	response := NewMockResponse()
	err := errors.New("test")

	Error(response, err)

	if response.Status != http.StatusInternalServerError {
		t.Errorf("response.Status is %v, expected %v (http.StatusInternalServerError)", response.Status, http.StatusNoContent)
	}

	if len(response.Data) == 0 {
		t.Errorf("len(response.Data) is %v, expected non-zero", len(response.Data))
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

	response := NewMockResponse()

	http_request := &http.Request{
		Method: request_method,
		URL:    request_url,
		Header: http.Header{},
	}

	request := NewRequest(http_request, request_path, module)

	Redirect(response, request, redirect_path)

	//response.Status has to be StatusSeeOther
	if response.Status != http.StatusSeeOther {
		t.Errorf("response.Status is %v, expected %v (http.StatusSeeOther)", response.Status, http.StatusSeeOther)
	}

	location := response.Headers.Get("Location")

	//the redirect must be relative to the module mount point
	if !strings.HasPrefix(location, module.MountPoint) {
		t.Errorf("Redirect Location header value is %v, expected this path to be prefixed with the module's mount point, %v", location, module.MountPoint)
	}

	//the redirect location must contain the redirect path
	if strings.Index(location, redirect_path) < len(module.MountPoint) {
		t.Errorf("Redirect Location header value is %v, expected this path to contain the redirect path %v", location, redirect_path)
	}

	//verify that X-Survana-Redirect isn't set
	_, ok := response.Headers["X-Survana-Redirect"]

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

	response := NewMockResponse()

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
	x_survana_redirect_header, ok := response.Headers["X-Survana-Redirect"]

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

	response := NewMockResponse()

	http_request := &http.Request{
		Method: request_method,
		URL:    request_url,
		Header: http.Header{},
	}

	request := NewRequest(http_request, request_path, module)

	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	XHRRedirect(response, request, redirect_path)

	//verify that X-Survana-Redirect is set
	x_survana_redirect_header, ok := response.Headers["X-Survana-Redirect"]

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
	response := NewMockResponse()
	expected_response := &JSONResponse{
		Success: true,
		Message: "This is a test",
	}

	JSONResult(response, expected_response.Success, expected_response.Message)

	if len(response.Data) == 0 {
		t.Errorf("len(response.Data) is %v, expected non-zero", len(response.Data))
	}

	//unmarshal the response
	actual_response := JSONResponse{}
	err := json.Unmarshal(response.Data, &actual_response)
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
	response := NewMockResponse()

	BadRequest(response)

	if response.Status != http.StatusBadRequest {
		t.Errorf("response.Status is %v, expected %v (http.StatusBadRequest)", response.Status, http.StatusBadRequest)
	}

	if len(response.Data) == 0 {
		t.Errorf("len(response.Data) is %v, expected non-zero", len(response.Data))
	}
}

func TestUnauthorized(t *testing.T) {
	response := NewMockResponse()
	err := errors.New("Access is denied")

	Unauthorized(response, err)

	if response.Status != http.StatusUnauthorized {
		t.Errorf("response.Status is %v, expected %v (http.StatusUnauthorized)", response.Status, http.StatusUnauthorized)
	}

	if len(response.Data) == 0 {
		t.Errorf("len(response.Data) is %v, expected non-zero", len(response.Data))
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

package perfect

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type routeTestCase struct {
	RoutePath, RequestPath string
}

var (
	routeTests = []routeTestCase{
		{"/", "/"},
		{"/test", "/test"},
		{"/test/test2", "/test/test2"},
		{"/:study", "/TEST_STUDY123"},
		{"/study/:study", "/study/TEST_STUDY_1234"},
		{"/study/:study/participants", "/study/A/participants"},
		{"/study/:study/form/:form", "/study/ABCD/form/ABCDEFG1234"},
		{"/study/:study/form/:form/subjects", "/study/abcd/form/foo/subjects"},
		{"/study/:study/form/:form/subjects/:subject", "/study/abcd/form/foo/subjects/ABC123"},
		{"/study/:study/form/:form/subjects/:subject/info", "/study/abcd/form/foo/subjects/abcd123/info"},
		{"/id/:id/id/:id/id/:id", "/id/1/id/2/id/3"},
		{"/id/:id/:id/:id", "/id/1/2/3"},
		{"/id/:id/:id/:id", "/id/1/2/3?id=4"},
		{"/:id/:name", "/abc/def"},
	}
)

func TestPrettyMux(t *testing.T) {
	_ = Module{
		Mux: NewPrettyMux(),
	}
}

func TestPrettyMux_FindHandler(t *testing.T) {
	var (
		called bool
	)

	m := NewPrettyMux()

	expected_handler := func(w http.ResponseWriter, r *Request) {
		called = true
	}

	//register all the routes first
	for _, test_case := range routeTests {
		m.Get(test_case.RoutePath, expected_handler)
	}

	for i, test_case := range routeTests {
		called = false
		/* Build a Request object */
		request_url, err := url.Parse("http://localhost" + test_case.RequestPath)
		if err != nil {
			t.Errorf("err = %v", err)
		}

		http_request := &http.Request{
			Method: "GET",
			URL:    request_url,
			Header: http.Header{},
		}

		v := strings.SplitN(test_case.RequestPath, "?", 2)
		request_path := v[0]

		request := NewRequest(http_request, request_path, &Module{})

		actual_handler := m.FindHandler(request)

		if actual_handler == nil {
			t.Fatalf("#%v: Actual handler is %v, expected non-nil", i, actual_handler)
		}

		response := httptest.NewRecorder()
		actual_handler(response, request)

		if !called {
			t.Fatalf("#%v: Actual handler is %v, expected handler is %v\n", i, actual_handler, expected_handler)
		}
	}
}

package perfect

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModuleMux(t *testing.T) {
	m := &Module{Mux: NewHTTPMux()}
	Modules.Mount(m, "/")
	server := httptest.NewServer(Modules)
	defer server.Close()

	response, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	t.Logf("body = %s", body)
}

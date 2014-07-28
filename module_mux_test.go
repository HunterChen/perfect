package perfect

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModuleMux(t *testing.T) {
	m := &Module{Mux: NewMux()}
	Modules.Mount(m, "/")
	server := httptest.NewServer(Modules)
	defer server.Close()

	response, err := http.Get(server.URL)
	if err != nil {
		t.Fatal("err = %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("body = %s", body)
}

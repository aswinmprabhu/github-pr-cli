package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Request(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	)
	defer ts.Close()
	testServerURL := ts.URL
	err := Publish(nsqdUrl, "hello")
	if err == nil {
		t.Errorf("Publish() didn’t return an error")
	}
}

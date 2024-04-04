package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/pes2324q2-gei-upc/ppf-chat-engine/api"
)

func TestRootHandler(t *testing.T) {
	// Create a new request to the root endpoint
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the rootHandler function
	api.RootHandler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	ts := newTestApplication(t).newTestServer(t)
	defer ts.Close()

	code, _, body := ts.get(t, "/")
	if code != http.StatusOK {
		t.Errorf("got: %d; wanted: %d", code, http.StatusOK)
	}

	expect := "Hello, world!"
	if !strings.Contains(body, expect) {
		t.Errorf("got: %s; expected to contain: %s", body, expect)
	}
}

package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestPingHandler(t *testing.T) {
	ts := newTestApplication(t).newTestServer(t)
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("got: %d; wanted: %d", code, http.StatusOK)
	}

	expect := "OK"
	if !strings.Contains(body, expect) {
		t.Errorf("got: %s; expected to contain: %s", body, expect)
	}
}

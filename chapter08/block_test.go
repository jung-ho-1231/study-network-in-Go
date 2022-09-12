package chapter08

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
	select {}
}

func TestBlockInfinitely(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
	_, _ = http.Get(ts.URL)
	t.Fatal("client not indefinitely block")
}

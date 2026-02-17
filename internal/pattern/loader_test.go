package pattern

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestShouldFailWhenResponseSizeExceedsLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(strings.Repeat("a", 20)))
	}))
	defer server.Close()

	loader := NewHTTPWikiLoader(1*time.Second, 10)
	_, err := loader.Load(server.URL)
	if err == nil {
		t.Fatalf("expected size limit error")
	}
}

func TestShouldFailWhenHTTPStatusIsNotOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	loader := NewHTTPWikiLoader(1*time.Second, 1024)
	_, err := loader.Load(server.URL)
	if err == nil {
		t.Fatalf("expected non-200 status error")
	}
}

func TestShouldFailWhenContentTypeIsNotText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	loader := NewHTTPWikiLoader(1*time.Second, 1024)
	_, err := loader.Load(server.URL)
	if err == nil {
		t.Fatalf("expected invalid content-type error")
	}
}

func TestShouldFailWhenRequestTimesOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("x = 1, y = 1\no!"))
	}))
	defer server.Close()

	loader := NewHTTPWikiLoader(10*time.Millisecond, 1024)
	_, err := loader.Load(server.URL)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
}

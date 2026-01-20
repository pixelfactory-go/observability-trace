package trace_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.pixelfactory.io/pkg/observability/trace"
)

func TestHTTPHandler(t *testing.T) {
	t.Parallel()

	t.Run("wraps handler correctly", func(t *testing.T) {
		t.Parallel()

		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		})

		wrappedHandler := trace.HTTPHandler(handler, "test-handler")

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		if !called {
			t.Error("original handler was not called")
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Body.String() != "OK" {
			t.Errorf("expected body %q, got %q", "OK", rec.Body.String())
		}
	})

	t.Run("preserves handler behavior", func(t *testing.T) {
		t.Parallel()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Custom-Header", "value")
			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write([]byte("Created")); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		})

		wrappedHandler := trace.HTTPHandler(handler, "create-handler")

		req := httptest.NewRequest(http.MethodPost, "/create", nil)
		rec := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
		}

		if rec.Header().Get("X-Custom-Header") != "value" {
			t.Error("custom header was not preserved")
		}

		if rec.Body.String() != "Created" {
			t.Errorf("expected body %q, got %q", "Created", rec.Body.String())
		}
	})
}

func TestHTTPHandlerFunc(t *testing.T) {
	t.Parallel()

	t.Run("wraps handler func correctly", func(t *testing.T) {
		t.Parallel()

		called := false
		handlerFunc := func(w http.ResponseWriter, _ *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}

		wrappedHandlerFunc := trace.HTTPHandlerFunc(handlerFunc, "test-handler-func")

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		wrappedHandlerFunc(rec, req)

		if !called {
			t.Error("original handler func was not called")
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Body.String() != "OK" {
			t.Errorf("expected body %q, got %q", "OK", rec.Body.String())
		}
	})
}

func TestHTTPClientTransporter(t *testing.T) {
	t.Parallel()

	t.Run("wraps round tripper", func(t *testing.T) {
		t.Parallel()

		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		// Create client with traced transport
		client := &http.Client{
			Transport: trace.HTTPClientTransporter(http.DefaultTransport),
		}

		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("preserves transport behavior", func(t *testing.T) {
		t.Parallel()

		// Create a test server that returns specific headers
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Server-Header", "test-value")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// Create client with traced transport
		client := &http.Client{
			Transport: trace.HTTPClientTransporter(http.DefaultTransport),
		}

		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.Header.Get("X-Server-Header") != "test-value" {
			t.Error("server header was not preserved")
		}
	})
}

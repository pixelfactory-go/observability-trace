package trace_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace/noop"

	trace "go.pixelfactory.io/pkg/observability/trace"
)

// FuzzAddSpanTags tests the AddSpanTags function with arbitrary string inputs.
func FuzzAddSpanTags(f *testing.F) {
	f.Add("key1", "value1")
	f.Add("", "")
	f.Add("user.id", "12345")
	f.Add("http.method", "GET")

	f.Fuzz(func(t *testing.T, key, value string) {
		_, span := noop.NewTracerProvider().Tracer("test").Start(context.Background(), "test")
		defer span.End()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("AddSpanTags panicked with key=%q, value=%q: %v", key, value, r)
			}
		}()

		tags := map[string]string{key: value}
		trace.AddSpanTags(span, tags)
	})
}

// FuzzAddSpanEvents tests the AddSpanEvents function with arbitrary string inputs.
func FuzzAddSpanEvents(f *testing.F) {
	f.Add("event1", "key1", "value1")
	f.Add("", "", "")
	f.Add("cache.hit", "cache.key", "user:123")
	f.Add("database.query", "query", "SELECT id, name FROM users")

	f.Fuzz(func(t *testing.T, name, key, value string) {
		_, span := noop.NewTracerProvider().Tracer("test").Start(context.Background(), "test")
		defer span.End()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("AddSpanEvents panicked with name=%q, key=%q, value=%q: %v", name, key, value, r)
			}
		}()

		events := map[string]string{key: value}
		trace.AddSpanEvents(span, name, events)
	})
}

// FuzzFailSpan tests the FailSpan function with arbitrary message strings.
func FuzzFailSpan(f *testing.F) {
	f.Add("error occurred")
	f.Add("")
	f.Add("database connection failed")
	f.Add("invalid input provided")

	f.Fuzz(func(t *testing.T, msg string) {
		_, span := noop.NewTracerProvider().Tracer("test").Start(context.Background(), "test")
		defer span.End()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("FailSpan panicked with msg=%q: %v", msg, r)
			}
		}()

		trace.FailSpan(span, msg)
	})
}

// FuzzNewSpan tests the NewSpan function with arbitrary span names.
func FuzzNewSpan(f *testing.F) {
	f.Add("operation")
	f.Add("")
	f.Add("http.request")
	f.Add("database.query")

	f.Fuzz(func(t *testing.T, name string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NewSpan panicked with name=%q: %v", name, r)
			}
		}()

		ctx, span := trace.NewSpan(context.Background(), name, nil)
		if ctx == nil {
			t.Error("NewSpan returned nil context")
		}
		if span == nil {
			t.Error("NewSpan returned nil span")
		}
		span.End()
	})
}

// FuzzWithServiceName tests the WithServiceName option with arbitrary service names.
func FuzzWithServiceName(f *testing.F) {
	f.Add("my-service")
	f.Add("")
	f.Add("api-gateway")
	f.Add("user-service")

	f.Fuzz(func(t *testing.T, name string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("WithServiceName panicked with name=%q: %v", name, r)
			}
		}()

		var c trace.Config
		opt := trace.WithServiceName(name)
		opt(&c)

		if c.ServiceName != name {
			t.Errorf("Expected ServiceName=%q, got %q", name, c.ServiceName)
		}
	})
}

// FuzzWithServiceVersion tests the WithServiceVersion option with arbitrary versions.
func FuzzWithServiceVersion(f *testing.F) {
	f.Add("1.0.0")
	f.Add("")
	f.Add("v2.3.4-beta")
	f.Add("latest")

	f.Fuzz(func(t *testing.T, version string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("WithServiceVersion panicked with version=%q: %v", version, r)
			}
		}()

		var c trace.Config
		opt := trace.WithServiceVersion(version)
		opt(&c)

		if c.ServiceVersion != version {
			t.Errorf("Expected ServiceVersion=%q, got %q", version, c.ServiceVersion)
		}
	})
}

// FuzzWithSpanExporterEndpoint tests the WithSpanExporterEndpoint option with arbitrary URLs.
func FuzzWithSpanExporterEndpoint(f *testing.F) {
	f.Add("http://localhost:4317")
	f.Add("")
	f.Add("https://otel-collector:4318")
	f.Add("grpc://collector.example.com:4317")

	f.Fuzz(func(t *testing.T, url string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("WithSpanExporterEndpoint panicked with url=%q: %v", url, r)
			}
		}()

		var c trace.Config
		opt := trace.WithSpanExporterEndpoint(url)
		opt(&c)

		if c.SpanExporterEndpoint != url {
			t.Errorf("Expected SpanExporterEndpoint=%q, got %q", url, c.SpanExporterEndpoint)
		}
	})
}

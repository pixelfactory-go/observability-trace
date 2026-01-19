package trace_test

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.pixelfactory.io/pkg/observability/trace"
)

func TestWithTraceEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{
			name:     "enable tracing",
			enabled:  true,
			expected: true,
		},
		{
			name:     "disable tracing",
			enabled:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg trace.Config
			opt := trace.WithTraceEnabled(tt.enabled)
			opt(&cfg)

			if cfg.TraceEnabled != tt.expected {
				t.Errorf("expected TraceEnabled=%v, got %v", tt.expected, cfg.TraceEnabled)
			}
		})
	}
}

func TestWithSpanExporterEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
	}{
		{
			name:     "localhost endpoint",
			endpoint: "localhost:4317",
		},
		{
			name:     "remote endpoint",
			endpoint: "collector.example.com:4317",
		},
		{
			name:     "http endpoint",
			endpoint: "http://localhost:4317",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg trace.Config
			opt := trace.WithSpanExporterEndpoint(tt.endpoint)
			opt(&cfg)

			if cfg.SpanExporterEndpoint != tt.endpoint {
				t.Errorf("expected SpanExporterEndpoint=%q, got %q", tt.endpoint, cfg.SpanExporterEndpoint)
			}
		})
	}
}

func TestWithServiceName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		serviceName string
	}{
		{
			name:        "simple service name",
			serviceName: "my-service",
		},
		{
			name:        "empty service name",
			serviceName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg trace.Config
			opt := trace.WithServiceName(tt.serviceName)
			opt(&cfg)

			if cfg.ServiceName != tt.serviceName {
				t.Errorf("expected ServiceName=%q, got %q", tt.serviceName, cfg.ServiceName)
			}
		})
	}
}

func TestWithServiceVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "semver version",
			version: "1.2.3",
		},
		{
			name:    "custom version",
			version: "v2024.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg trace.Config
			opt := trace.WithServiceVersion(tt.version)
			opt(&cfg)

			if cfg.ServiceVersion != tt.version {
				t.Errorf("expected ServiceVersion=%q, got %q", tt.version, cfg.ServiceVersion)
			}
		})
	}
}

func TestWithResourceAttributes(t *testing.T) {
	t.Parallel()

	attributes := map[string]string{
		"environment": "production",
		"region":      "us-east-1",
	}

	var cfg trace.Config
	opt := trace.WithResourceAttributes(attributes)
	opt(&cfg)

	if len(cfg.ResourceAttributes) != len(attributes) {
		t.Errorf("expected %d attributes, got %d", len(attributes), len(cfg.ResourceAttributes))
	}

	for key, expectedValue := range attributes {
		if actualValue, ok := cfg.ResourceAttributes[key]; !ok {
			t.Errorf("attribute %q not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("attribute %q: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestWithPropagators(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		propagators []string
	}{
		{
			name:        "single propagator",
			propagators: []string{"b3"},
		},
		{
			name:        "multiple propagators",
			propagators: []string{"b3", "tracecontext", "baggage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg trace.Config
			opt := trace.WithPropagators(tt.propagators)
			opt(&cfg)

			if len(cfg.Propagators) != len(tt.propagators) {
				t.Errorf("expected %d propagators, got %d", len(tt.propagators), len(cfg.Propagators))
			}

			for i, expected := range tt.propagators {
				if cfg.Propagators[i] != expected {
					t.Errorf("propagator[%d]: expected %q, got %q", i, expected, cfg.Propagators[i])
				}
			}
		})
	}
}

func TestWithHeaders(t *testing.T) {
	t.Parallel()

	t.Run("add headers to empty config", func(t *testing.T) {
		t.Parallel()

		headers := map[string]string{
			"api-key":      "secret",
			"x-custom-hdr": "value",
		}

		var cfg trace.Config
		opt := trace.WithHeaders(headers)
		opt(&cfg)

		if len(cfg.Headers) != len(headers) {
			t.Errorf("expected %d headers, got %d", len(headers), len(cfg.Headers))
		}

		for key, expectedValue := range headers {
			if actualValue, ok := cfg.Headers[key]; !ok {
				t.Errorf("header %q not found", key)
			} else if actualValue != expectedValue {
				t.Errorf("header %q: expected %q, got %q", key, expectedValue, actualValue)
			}
		}
	})

	t.Run("add headers to existing config", func(t *testing.T) {
		t.Parallel()

		cfg := trace.Config{
			Headers: map[string]string{
				"existing": "header",
			},
		}

		newHeaders := map[string]string{
			"new": "value",
		}

		opt := trace.WithHeaders(newHeaders)
		opt(&cfg)

		if len(cfg.Headers) != 2 {
			t.Errorf("expected 2 headers, got %d", len(cfg.Headers))
		}

		if cfg.Headers["existing"] != "header" {
			t.Error("existing header was lost")
		}

		if cfg.Headers["new"] != "value" {
			t.Error("new header was not added")
		}
	})
}

func TestWithSpanExporterInsecure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		insecure bool
	}{
		{
			name:     "insecure enabled",
			insecure: true,
		},
		{
			name:     "insecure disabled",
			insecure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg trace.Config
			opt := trace.WithSpanExporterInsecure(tt.insecure)
			opt(&cfg)

			if cfg.SpanExporterEndpointInsecure != tt.insecure {
				t.Errorf("expected SpanExporterEndpointInsecure=%v, got %v", tt.insecure, cfg.SpanExporterEndpointInsecure)
			}
		})
	}
}

func TestWithTraceExporter(t *testing.T) {
	t.Parallel()

	// Create a mock exporter
	exporter := &mockSpanExporter{}

	var cfg trace.Config
	opt := trace.WithTraceExporter(exporter)
	opt(&cfg)

	if cfg.TraceExporter != exporter {
		t.Error("TraceExporter was not set correctly")
	}
}

type mockSpanExporter struct{}

func (m *mockSpanExporter) ExportSpans(_ context.Context, _ []sdktrace.ReadOnlySpan) error {
	return nil
}

func (m *mockSpanExporter) Shutdown(_ context.Context) error {
	return nil
}

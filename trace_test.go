package trace_test

import (
	"testing"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.pixelfactory.io/pkg/observability/trace"
)

func TestNewProvider(t *testing.T) {
	t.Parallel()

	t.Run("create provider with tracing disabled", func(t *testing.T) {
		t.Parallel()

		provider, err := trace.NewProvider(
			trace.WithTraceEnabled(false),
			trace.WithServiceName("test-service"),
		)
		if err != nil {
			t.Fatalf("NewProvider failed: %v", err)
		}

		if provider == nil {
			t.Fatal("expected provider, got nil")
		}

		// Shutdown should work even with tracing disabled
		if err := provider.Shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})

	t.Run("create provider with tracing enabled", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create exporter: %v", err)
		}

		provider, err := trace.NewProvider(
			trace.WithTraceEnabled(true),
			trace.WithServiceName("test-service"),
			trace.WithServiceVersion("1.0.0"),
			trace.WithTraceExporter(exp),
			trace.WithPropagators([]string{"b3"}),
		)
		if err != nil {
			t.Fatalf("NewProvider failed: %v", err)
		}

		if provider == nil {
			t.Fatal("expected provider, got nil")
		}

		// Cleanup
		if err := provider.Shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})

	t.Run("create provider with all options", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create exporter: %v", err)
		}

		provider, err := trace.NewProvider(
			trace.WithTraceEnabled(true),
			trace.WithServiceName("comprehensive-service"),
			trace.WithServiceVersion("2.0.0"),
			trace.WithSpanExporterEndpoint("localhost:4317"),
			trace.WithSpanExporterInsecure(true),
			trace.WithTraceExporter(exp),
			trace.WithPropagators([]string{"b3", "tracecontext", "baggage"}),
			trace.WithHeaders(map[string]string{"api-key": "secret"}),
			trace.WithResourceAttributes(map[string]string{
				"environment": "test",
				"region":      "us-east-1",
			}),
		)
		if err != nil {
			t.Fatalf("NewProvider failed: %v", err)
		}

		if provider == nil {
			t.Fatal("expected provider, got nil")
		}

		// Cleanup
		if err := provider.Shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})

	t.Run("create provider with invalid propagators", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create exporter: %v", err)
		}

		_, err = trace.NewProvider(
			trace.WithTraceEnabled(true),
			trace.WithServiceName("test-service"),
			trace.WithTraceExporter(exp),
			trace.WithPropagators([]string{"invalid-propagator"}),
		)

		if err == nil {
			t.Error("expected error for invalid propagators, got nil")
		}
	})

	t.Run("provider shutdown can be called multiple times", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create exporter: %v", err)
		}

		provider, err := trace.NewProvider(
			trace.WithTraceEnabled(true),
			trace.WithServiceName("test-service"),
			trace.WithTraceExporter(exp),
			trace.WithPropagators([]string{"b3"}),
		)
		if err != nil {
			t.Fatalf("NewProvider failed: %v", err)
		}

		// Call shutdown multiple times
		if err := provider.Shutdown(); err != nil {
			t.Errorf("first shutdown failed: %v", err)
		}

		if err := provider.Shutdown(); err != nil {
			t.Errorf("second shutdown failed: %v", err)
		}
	})
}

func TestProviderWithCustomExporter(t *testing.T) {
	t.Parallel()

	t.Run("use custom exporter", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
		if err != nil {
			t.Fatalf("failed to create exporter: %v", err)
		}

		provider, err := trace.NewProvider(
			trace.WithTraceEnabled(true),
			trace.WithServiceName("custom-exporter-service"),
			trace.WithTraceExporter(exp),
			trace.WithPropagators([]string{"b3"}),
		)
		if err != nil {
			t.Fatalf("NewProvider failed: %v", err)
		}

		if provider == nil {
			t.Fatal("expected provider, got nil")
		}

		// Cleanup
		if err := provider.Shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})
}

func TestProviderHeadersMerge(t *testing.T) {
	t.Parallel()

	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	provider, err := trace.NewProvider(
		trace.WithTraceEnabled(true),
		trace.WithServiceName("headers-service"),
		trace.WithTraceExporter(exp),
		trace.WithPropagators([]string{"b3"}),
		trace.WithHeaders(map[string]string{"header1": "value1"}),
		trace.WithHeaders(map[string]string{"header2": "value2"}),
	)
	if err != nil {
		t.Fatalf("NewProvider failed: %v", err)
	}

	if provider == nil {
		t.Fatal("expected provider, got nil")
	}

	// Cleanup
	if err := provider.Shutdown(); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

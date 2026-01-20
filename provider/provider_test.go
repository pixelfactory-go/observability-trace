package provider_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"go.pixelfactory.io/pkg/observability/trace/provider"
)

func createTestExporter(t *testing.T) tracesdk.SpanExporter {
	t.Helper()
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Fatalf("failed to create stdout exporter: %v", err)
	}
	return exp
}

func createTestResource(t *testing.T) *resource.Resource {
	t.Helper()
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-service"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	return res
}

func testInitProviderSuccess(t *testing.T, cfg provider.Config) {
	t.Helper()
	shutdown, err := provider.InitProvider(cfg)
	if err != nil {
		t.Fatalf("InitProvider failed: %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected shutdown function, got nil")
	}

	// Cleanup
	if shutdownErr := shutdown(); shutdownErr != nil {
		t.Errorf("shutdown failed: %v", shutdownErr)
	}
}

func TestInitProvider(t *testing.T) {
	t.Parallel()

	t.Run("initialize with valid config", func(t *testing.T) {
		t.Parallel()

		cfg := provider.Config{
			Endpoint:      "localhost:4317",
			Insecure:      true,
			Headers:       map[string]string{},
			Resource:      createTestResource(t),
			TraceExporter: createTestExporter(t),
			Propagators:   []string{"b3"},
		}

		shutdown, err := provider.InitProvider(cfg)
		if err != nil {
			t.Fatalf("InitProvider failed: %v", err)
		}

		if shutdown == nil {
			t.Fatal("expected shutdown function, got nil")
		}

		// Verify tracer provider was set
		if otel.GetTracerProvider() == nil {
			t.Error("expected tracer provider to be set")
		}

		// Cleanup
		if shutdownErr := shutdown(); shutdownErr != nil {
			t.Errorf("shutdown failed: %v", shutdownErr)
		}
	})

	t.Run("initialize with invalid propagators", func(t *testing.T) {
		t.Parallel()

		cfg := provider.Config{
			Endpoint:    "localhost:4317",
			Insecure:    true,
			Headers:     map[string]string{},
			Resource:    createTestResource(t),
			Propagators: []string{"invalid-propagator"},
		}

		_, err := provider.InitProvider(cfg)
		if err == nil {
			t.Error("expected error for invalid propagators, got nil")
		}

		expectedErrMsg := "invalid configuration: unsupported propagators. Supported options: b3,baggage,tracecontext,ottrace"
		if err.Error() != expectedErrMsg {
			t.Errorf("expected error %q, got %q", expectedErrMsg, err.Error())
		}
	})

	t.Run("initialize with empty propagators", func(t *testing.T) {
		t.Parallel()

		cfg := provider.Config{
			Endpoint:    "localhost:4317",
			Insecure:    true,
			Headers:     map[string]string{},
			Resource:    createTestResource(t),
			Propagators: []string{},
		}

		_, err := provider.InitProvider(cfg)
		if err == nil {
			t.Error("expected error for empty propagators, got nil")
		}
	})

	t.Run("initialize with multiple valid propagators", func(t *testing.T) {
		t.Parallel()

		cfg := provider.Config{
			Endpoint:      "localhost:4317",
			Insecure:      true,
			Headers:       map[string]string{},
			Resource:      createTestResource(t),
			TraceExporter: createTestExporter(t),
			Propagators:   []string{"b3", "tracecontext", "baggage", "ottrace"},
		}

		testInitProviderSuccess(t, cfg)
	})

	t.Run("initialize with custom headers", func(t *testing.T) {
		t.Parallel()

		cfg := provider.Config{
			Endpoint:      "localhost:4317",
			Insecure:      true,
			Headers:       map[string]string{"api-key": "secret", "x-custom": "value"},
			Resource:      createTestResource(t),
			TraceExporter: createTestExporter(t),
			Propagators:   []string{"b3"},
		}

		testInitProviderSuccess(t, cfg)
	})
}

func TestShutdownFunc(t *testing.T) {
	t.Parallel()

	cfg := provider.Config{
		Endpoint:      "localhost:4317",
		Insecure:      true,
		Headers:       map[string]string{},
		Resource:      createTestResource(t),
		TraceExporter: createTestExporter(t),
		Propagators:   []string{"b3"},
	}

	shutdown, err := provider.InitProvider(cfg)
	if err != nil {
		t.Fatalf("InitProvider failed: %v", err)
	}

	// Test that shutdown can be called
	if shutdownErr := shutdown(); shutdownErr != nil {
		t.Errorf("shutdown failed: %v", shutdownErr)
	}

	// Test that shutdown can be called multiple times
	if shutdownErr := shutdown(); shutdownErr != nil {
		t.Errorf("second shutdown failed: %v", shutdownErr)
	}
}

package provider_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"go.pixelfactory.io/pkg/observability/trace/provider"
)

func TestInitProvider(t *testing.T) {
	t.Parallel()

	t.Run("initialize with valid config", func(t *testing.T) {
		t.Parallel()

		// Create a stdout exporter for testing
		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create stdout exporter: %v", err)
		}

		res, err := resource.New(
			context.Background(),
			resource.WithAttributes(
				semconv.ServiceName("test-service"),
			),
		)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		cfg := provider.Config{
			Endpoint:      "localhost:4317",
			Insecure:      true,
			Headers:       map[string]string{},
			Resource:      res,
			TraceExporter: exp,
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
		if err := shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})

	t.Run("initialize with invalid propagators", func(t *testing.T) {
		t.Parallel()

		res, err := resource.New(
			context.Background(),
			resource.WithAttributes(
				semconv.ServiceName("test-service"),
			),
		)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		cfg := provider.Config{
			Endpoint:    "localhost:4317",
			Insecure:    true,
			Headers:     map[string]string{},
			Resource:    res,
			Propagators: []string{"invalid-propagator"},
		}

		_, err = provider.InitProvider(cfg)
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

		res, err := resource.New(
			context.Background(),
			resource.WithAttributes(
				semconv.ServiceName("test-service"),
			),
		)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		cfg := provider.Config{
			Endpoint:    "localhost:4317",
			Insecure:    true,
			Headers:     map[string]string{},
			Resource:    res,
			Propagators: []string{},
		}

		_, err = provider.InitProvider(cfg)
		if err == nil {
			t.Error("expected error for empty propagators, got nil")
		}
	})

	t.Run("initialize with multiple valid propagators", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create stdout exporter: %v", err)
		}

		res, err := resource.New(
			context.Background(),
			resource.WithAttributes(
				semconv.ServiceName("test-service"),
			),
		)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		cfg := provider.Config{
			Endpoint:      "localhost:4317",
			Insecure:      true,
			Headers:       map[string]string{},
			Resource:      res,
			TraceExporter: exp,
			Propagators:   []string{"b3", "tracecontext", "baggage", "ottrace"},
		}

		shutdown, err := provider.InitProvider(cfg)
		if err != nil {
			t.Fatalf("InitProvider failed: %v", err)
		}

		if shutdown == nil {
			t.Fatal("expected shutdown function, got nil")
		}

		// Cleanup
		if err := shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})

	t.Run("initialize with custom headers", func(t *testing.T) {
		t.Parallel()

		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			t.Fatalf("failed to create stdout exporter: %v", err)
		}

		res, err := resource.New(
			context.Background(),
			resource.WithAttributes(
				semconv.ServiceName("test-service"),
			),
		)
		if err != nil {
			t.Fatalf("failed to create resource: %v", err)
		}

		cfg := provider.Config{
			Endpoint:      "localhost:4317",
			Insecure:      true,
			Headers:       map[string]string{"api-key": "secret", "x-custom": "value"},
			Resource:      res,
			TraceExporter: exp,
			Propagators:   []string{"b3"},
		}

		shutdown, err := provider.InitProvider(cfg)
		if err != nil {
			t.Fatalf("InitProvider failed: %v", err)
		}

		if shutdown == nil {
			t.Fatal("expected shutdown function, got nil")
		}

		// Cleanup
		if err := shutdown(); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})
}

func TestShutdownFunc(t *testing.T) {
	t.Parallel()

	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Fatalf("failed to create stdout exporter: %v", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-service"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	cfg := provider.Config{
		Endpoint:      "localhost:4317",
		Insecure:      true,
		Headers:       map[string]string{},
		Resource:      res,
		TraceExporter: exp,
		Propagators:   []string{"b3"},
	}

	shutdown, err := provider.InitProvider(cfg)
	if err != nil {
		t.Fatalf("InitProvider failed: %v", err)
	}

	// Test that shutdown can be called
	if err := shutdown(); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}

	// Test that shutdown can be called multiple times
	if err := shutdown(); err != nil {
		t.Errorf("second shutdown failed: %v", err)
	}
}

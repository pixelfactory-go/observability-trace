package trace_test

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"

	"go.pixelfactory.io/pkg/observability/trace"
)

func setupTestTracer() (*tracetest.SpanRecorder, func()) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	otel.SetTracerProvider(tp)

	cleanup := func() {
		_ = tp.Shutdown(context.Background())
	}

	return sr, cleanup
}

func TestNewSpan(t *testing.T) {
	t.Parallel()

	t.Run("create span without customizer", func(t *testing.T) {
		t.Parallel()
		sr, cleanup := setupTestTracer()
		defer cleanup()

		ctx := context.Background()
		spanName := "test-span"

		newCtx, span := trace.NewSpan(ctx, spanName, nil)
		span.End()

		if newCtx == ctx {
			t.Error("expected new context to be different from original")
		}

		spans := sr.Ended()
		if len(spans) != 1 {
			t.Fatalf("expected 1 span, got %d", len(spans))
		}

		if spans[0].Name() != spanName {
			t.Errorf("expected span name %q, got %q", spanName, spans[0].Name())
		}
	})

	t.Run("create span with customizer", func(t *testing.T) {
		t.Parallel()
		sr, cleanup := setupTestTracer()
		defer cleanup()

		ctx := context.Background()
		spanName := "test-span-custom"

		customizer := &testSpanCustomiser{
			kind: oteltrace.SpanKindClient,
		}

		newCtx, span := trace.NewSpan(ctx, spanName, customizer)
		span.End()

		if newCtx == ctx {
			t.Error("expected new context to be different from original")
		}

		spans := sr.Ended()
		if len(spans) != 1 {
			t.Fatalf("expected 1 span, got %d", len(spans))
		}

		if spans[0].Name() != spanName {
			t.Errorf("expected span name %q, got %q", spanName, spans[0].Name())
		}

		if spans[0].SpanKind() != oteltrace.SpanKindClient {
			t.Errorf("expected span kind %v, got %v", oteltrace.SpanKindClient, spans[0].SpanKind())
		}
	})
}

func TestSpanFromContext(t *testing.T) {
	t.Parallel()

	_, cleanup := setupTestTracer()
	defer cleanup()

	ctx := context.Background()
	ctx, originalSpan := trace.NewSpan(ctx, "parent-span", nil)
	defer originalSpan.End()

	retrievedSpan := trace.SpanFromContext(ctx)

	if retrievedSpan.SpanContext().SpanID() != originalSpan.SpanContext().SpanID() {
		t.Error("retrieved span does not match original span")
	}
}

func TestAddSpanTags(t *testing.T) {
	t.Parallel()

	sr, cleanup := setupTestTracer()
	defer cleanup()

	ctx := context.Background()
	_, span := trace.NewSpan(ctx, "test-span", nil)

	tags := map[string]string{
		"user.id": "12345",
		"action":  "create",
		"region":  "us-east-1",
	}

	trace.AddSpanTags(span, tags)
	span.End()

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	attributes := spans[0].Attributes()
	if len(attributes) != len(tags) {
		t.Errorf("expected %d attributes, got %d", len(tags), len(attributes))
	}

	attrMap := make(map[string]string)
	for _, attr := range attributes {
		attrMap[string(attr.Key)] = attr.Value.AsString()
	}

	for key, expectedValue := range tags {
		if actualValue, ok := attrMap[key]; !ok {
			t.Errorf("attribute %q not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("attribute %q: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestAddSpanEvents(t *testing.T) {
	t.Parallel()

	sr, cleanup := setupTestTracer()
	defer cleanup()

	ctx := context.Background()
	_, span := trace.NewSpan(ctx, "test-span", nil)

	eventName := "processing-started"
	events := map[string]string{
		"timestamp": "2024-01-01T00:00:00Z",
		"step":      "initialization",
	}

	trace.AddSpanEvents(span, eventName, events)
	span.End()

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	spanEvents := spans[0].Events()
	if len(spanEvents) != 1 {
		t.Fatalf("expected 1 event, got %d", len(spanEvents))
	}

	if spanEvents[0].Name != eventName {
		t.Errorf("expected event name %q, got %q", eventName, spanEvents[0].Name)
	}

	if len(spanEvents[0].Attributes) != len(events) {
		t.Errorf("expected %d event attributes, got %d", len(events), len(spanEvents[0].Attributes))
	}
}

func TestAddSpanError(t *testing.T) {
	t.Parallel()

	sr, cleanup := setupTestTracer()
	defer cleanup()

	ctx := context.Background()
	_, span := trace.NewSpan(ctx, "test-span", nil)

	testErr := errors.New("test error")
	trace.AddSpanError(span, testErr)
	span.End()

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	spanEvents := spans[0].Events()
	if len(spanEvents) != 1 {
		t.Fatalf("expected 1 event (error), got %d", len(spanEvents))
	}

	if spanEvents[0].Name != "exception" {
		t.Errorf("expected event name %q, got %q", "exception", spanEvents[0].Name)
	}
}

func TestFailSpan(t *testing.T) {
	t.Parallel()

	sr, cleanup := setupTestTracer()
	defer cleanup()

	ctx := context.Background()
	_, span := trace.NewSpan(ctx, "test-span", nil)

	errorMessage := "operation failed"
	trace.FailSpan(span, errorMessage)
	span.End()

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	status := spans[0].Status()
	if status.Code != codes.Error {
		t.Errorf("expected status code %v, got %v", codes.Error, status.Code)
	}

	if status.Description != errorMessage {
		t.Errorf("expected status description %q, got %q", errorMessage, status.Description)
	}
}

// testSpanCustomiser implements SpanCustomiser for testing.
type testSpanCustomiser struct {
	kind oteltrace.SpanKind
}

func (t *testSpanCustomiser) Customise() []oteltrace.SpanStartOption {
	return []oteltrace.SpanStartOption{
		oteltrace.WithSpanKind(t.kind),
	}
}

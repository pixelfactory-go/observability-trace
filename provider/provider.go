package provider

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

type Config struct {
	Endpoint        string
	Insecure        bool
	Headers         map[string]string
	Resource        *resource.Resource
	TraceExporter   trace.SpanExporter
	ReportingPeriod string
	Propagators     []string
}

type ShutdownFunc func() error

type SetupFunc func(Config) (ShutdownFunc, error)

func InitProvider(c Config) (ShutdownFunc, error) {
	ctx := context.Background()
	var bsp trace.SpanProcessor

	traceExporter, err := newTraceExporter(c.Endpoint, c.Insecure, c.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to create span exporter: %w", err)
	}
	bsp = trace.NewBatchSpanProcessor(traceExporter)

	if c.TraceExporter != nil {
		bsp = trace.NewBatchSpanProcessor(c.TraceExporter)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(bsp),
		trace.WithResource(c.Resource),
	)

	if cfgErr := configurePropagators(c); cfgErr != nil {
		return nil, cfgErr
	}

	otel.SetTracerProvider(tracerProvider)

	return func() error {
		// Shutdown will flush any remaining spans and shut down the exporter.
		return tracerProvider.Shutdown(ctx)
	}, nil
}

func newTraceExporter(endpoint string, insecure bool, headers map[string]string) (*otlptrace.Exporter, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}
	return otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithHeaders(headers),
			otlptracegrpc.WithCompressor(gzip.Name),
		),
	)
}

// configurePropagators configures B3 propagation by default.
func configurePropagators(c Config) error {
	propagatorsMap := map[string]propagation.TextMapPropagator{
		"b3":           b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)),
		"baggage":      propagation.Baggage{},
		"tracecontext": propagation.TraceContext{},
		"ottrace":      ot.OT{},
	}
	var props []propagation.TextMapPropagator
	for _, key := range c.Propagators {
		prop := propagatorsMap[key]
		if prop != nil {
			props = append(props, prop)
		}
	}
	if len(props) == 0 {
		return errors.New(
			"invalid configuration: unsupported propagators. Supported options: b3,baggage,tracecontext,ottrace",
		)
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		props...,
	))
	return nil
}

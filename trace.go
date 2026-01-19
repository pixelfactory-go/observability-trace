package trace

import (
	"context"
	"os"

	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.pixelfactory.io/pkg/observability/trace/provider"
)

type Provider struct {
	config       Config
	ShutdownFunc provider.ShutdownFunc
}

func newConfig(opts ...Option) (Config, error) {
	var c Config
	err := envconfig.Process(context.Background(), &c)
	if err != nil {
		return Config{}, err
	}

	var defaultOpts []Option
	for _, opt := range append(defaultOpts, opts...) {
		opt(&c)
	}

	res, err := newResource(&c)
	if err != nil {
		return Config{}, err
	}
	c.Resource = res

	return c, nil
}

func newResource(c *Config) (*resource.Resource, error) {
	r := resource.Environment()

	hostnameSet := false
	for iter := r.Iter(); iter.Next(); {
		if iter.Attribute().Key == semconv.HostNameKey && len(iter.Attribute().Value.Emit()) > 0 {
			hostnameSet = true
		}
	}

	attributes := []attribute.KeyValue{
		semconv.TelemetrySDKNameKey.String("go.pixelfactory.io/pkg/observability/trace"),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersionKey.String(version),
	}

	if len(c.ServiceName) > 0 {
		attributes = append(attributes, semconv.ServiceNameKey.String(c.ServiceName))
	}

	if len(c.ServiceVersion) > 0 {
		attributes = append(attributes, semconv.ServiceVersionKey.String(c.ServiceVersion))
	}

	for key, value := range c.ResourceAttributes {
		if len(value) > 0 {
			if key == string(semconv.HostNameKey) {
				hostnameSet = true
			}
			attributes = append(attributes, semconv.HostNameKey.String(value))
		}
	}

	if !hostnameSet {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, semconv.HostNameKey.String(hostname))
	}

	attributes = append(r.Attributes(), attributes...)

	// These detectors can't actually fail, ignoring the error.
	r, _ = resource.New(
		context.Background(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(attributes...),
	)

	return r, nil
}

func setupTracing(c Config) (provider.ShutdownFunc, error) {
	if !c.TraceEnabled {
		return func() error { return nil }, nil
	}
	return provider.InitProvider(provider.Config{
		Endpoint:      c.SpanExporterEndpoint,
		Insecure:      c.SpanExporterEndpointInsecure,
		Headers:       c.Headers,
		Resource:      c.Resource,
		Propagators:   c.Propagators,
		TraceExporter: c.TraceExporter,
	})
}

// NewProvider returns a new `Provider` type.
func NewProvider(opts ...Option) (*Provider, error) {
	c, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	if c.Headers == nil {
		c.Headers = map[string]string{}
	}

	shutdown, err := setupTracing(c)
	if err != nil {
		return nil, err
	}

	p := &Provider{
		config:       c,
		ShutdownFunc: shutdown,
	}

	return p, nil
}

func (p Provider) Shutdown() error {
	return p.ShutdownFunc()
}

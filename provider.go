package trace

import (
	"context"
	"os"

	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.pixelfactory.io/pkg/observability/trace/pipelines"
)

type Provider struct {
	config        Config
	shutdownFuncs []func() error
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
		} else {
			attributes = append(attributes, semconv.HostNameKey.String(hostname))
		}
	}

	attributes = append(r.Attributes(), attributes...)

	// These detectors can't actually fail, ignoring the error.
	r, _ = resource.New(
		context.Background(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(attributes...),
	)

	// Note: There are new detectors we may wish to take advantage
	// of, now available in the default SDK (e.g., WithProcess(),
	// WithOSType(), ...).
	return r, nil
}

type setupFunc func(Config) (func() error, error)

func setupTracing(c Config) (func() error, error) {
	// if c.SpanExporterEndpoint == "" {
	// 	c.logger.Debugf("tracing is disabled by configuration: no endpoint set")
	// 	return nil, nil
	// }
	return pipelines.NewTracePipeline(pipelines.PipelineConfig{
		Endpoint:    c.SpanExporterEndpoint,
		Insecure:    c.SpanExporterEndpointInsecure,
		Headers:     c.Headers,
		Resource:    c.Resource,
		Propagators: c.Propagators,
	})
}

// New returns a new `Provider` type. It uses Jaeger exporter and globally sets
// the tracer provider as well as the global tracer for spans.
func NewProvider(opts ...Option) (Provider, error) {
	c, err := newConfig(opts...)
	if err != nil {
		return Provider{}, nil
	}

	if c.Headers == nil {
		c.Headers = map[string]string{}
	}

	p := Provider{
		config: c,
	}

	for _, setup := range []setupFunc{setupTracing} {
		shutdown, err := setup(c)
		if err != nil {
			return p, err
		}
		if shutdown != nil {
			p.shutdownFuncs = append(p.shutdownFuncs, shutdown)
		}
	}

	return p, nil
}

func (p Provider) Shutdown() {
	for _, shutdown := range p.shutdownFuncs {
		if err := shutdown(); err != nil {
			panic(err)
		}
	}
}

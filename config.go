package trace

import (
	"go.opentelemetry.io/otel/sdk/resource"
)

type Config struct {
	SpanExporterEndpoint           string            `env:"OTEL_EXPORTER_OTLP_SPAN_ENDPOINT"`
	SpanExporterEndpointInsecure   bool              `env:"OTEL_EXPORTER_OTLP_SPAN_INSECURE,default=false"`
	ServiceName                    string            `env:"LS_SERVICE_NAME"`
	ServiceVersion                 string            `env:"LS_SERVICE_VERSION,default=unknown"`
	Headers                        map[string]string `env:"OTEL_EXPORTER_OTLP_HEADERS"`
	MetricExporterEndpoint         string            `env:"OTEL_EXPORTER_OTLP_METRIC_ENDPOINT"`
	MetricExporterEndpointInsecure bool              `env:"OTEL_EXPORTER_OTLP_METRIC_INSECURE,default=false"`
	MetricsEnabled                 bool              `env:"LS_METRICS_ENABLED,default=true"`
	LogLevel                       string            `env:"OTEL_LOG_LEVEL,default=info"`
	Propagators                    []string          `env:"OTEL_PROPAGATORS,default=b3"`
	MetricReportingPeriod          string            `env:"OTEL_EXPORTER_OTLP_METRIC_PERIOD,default=30s"`
	ResourceAttributes             map[string]string
	Resource                       *resource.Resource
	Disabled                       bool
}

type Option func(*Config)

// WithSpanExporterEndpoint configures the endpoint for sending traces via OTLP
func WithSpanExporterEndpoint(url string) Option {
	return func(c *Config) {
		c.SpanExporterEndpoint = url
	}
}

// WithServiceName configures a "service.name" resource label
func WithServiceName(name string) Option {
	return func(c *Config) {
		c.ServiceName = name
	}
}

// WithServiceVersion configures a "service.version" resource label
func WithServiceVersion(version string) Option {
	return func(c *Config) {
		c.ServiceVersion = version
	}
}

// WithResourceAttributes configures attributes on the resource
func WithResourceAttributes(attributes map[string]string) Option {
	return func(c *Config) {
		c.ResourceAttributes = attributes
	}
}

// WithPropagators configures propagators
func WithPropagators(propagators []string) Option {
	return func(c *Config) {
		c.Propagators = propagators
	}
}

// WithHeaders configures OTLP/gRPC connection headers
func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// WithSpanExporterInsecure permits connecting to the
// trace endpoint without a certificate
func WithSpanExporterInsecure(insecure bool) Option {
	return func(c *Config) {
		c.SpanExporterEndpointInsecure = insecure
	}
}

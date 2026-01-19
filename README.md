# Trace

[![Go Reference](https://pkg.go.dev/badge/go.pixelfactory.io/pkg/observability/trace.svg)](https://pkg.go.dev/go.pixelfactory.io/pkg/observability/trace)
[![Go Report Card](https://goreportcard.com/badge/go.pixelfactory.io/pkg/observability/trace)](https://goreportcard.com/report/go.pixelfactory.io/pkg/observability/trace)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A batteries-included OpenTelemetry trace library for Go applications that simplifies distributed tracing setup and usage.

## Features

- **Simple Setup**: Initialize OpenTelemetry tracing with minimal configuration
- **Environment-based Configuration**: Configure via environment variables following OpenTelemetry standards
- **HTTP Instrumentation**: Built-in wrappers for HTTP servers and clients
- **Span Helpers**: Convenient functions for creating and managing spans
- **Multiple Propagators**: Support for B3, W3C TraceContext, Baggage, and OT propagators
- **OTLP Support**: Native gRPC export to OpenTelemetry collectors
- **Flexible Exporters**: Use OTLP or custom exporters (stdout, Jaeger, etc.)

## Installation

```bash
go get go.pixelfactory.io/pkg/observability/trace
```

## Usage

### Basic Setup

```go
package main

import (
    "log"
    "net/http"

    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.pixelfactory.io/pkg/observability/trace"
)

func main() {
    // Create console exporter for development
    exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
    if err != nil {
        log.Fatal(err)
    }

    // Initialize trace provider
    provider, err := trace.NewProvider(
        trace.WithTraceExporter(exp),
        trace.WithServiceName("my-service"),
        trace.WithServiceVersion("1.0.0"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer provider.Shutdown()

    // Your application code here
}
```

### HTTP Server Instrumentation

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "go.pixelfactory.io/pkg/observability/trace"
)

func main() {
    // Initialize provider (omitted for brevity)
    provider, _ := trace.NewProvider(
        trace.WithTraceEnabled(true),
        trace.WithServiceName("http-server"),
    )
    defer provider.Shutdown()

    // Define your handler
    helloHandler := func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(w, "Hello, World!\n")
    }

    // Wrap handler with tracing
    tracedHandler := trace.HTTPHandler(http.HandlerFunc(helloHandler), "hello-endpoint")

    http.Handle("/hello", tracedHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### HTTP Client Instrumentation

```go
package main

import (
    "context"
    "log"
    "net/http"

    "go.pixelfactory.io/pkg/observability/trace"
)

func main() {
    // Initialize provider
    provider, _ := trace.NewProvider(
        trace.WithTraceEnabled(true),
        trace.WithServiceName("http-client"),
    )
    defer provider.Shutdown()

    // Create traced HTTP client
    client := &http.Client{
        Transport: trace.HTTPClientTransporter(http.DefaultTransport),
    }

    req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
}
```

### Custom Spans

```go
package main

import (
    "context"

    "go.pixelfactory.io/pkg/observability/trace"
)

func processData(ctx context.Context) error {
    // Create a new span
    ctx, span := trace.NewSpan(ctx, "process-data", nil)
    defer span.End()

    // Add tags to the span
    trace.AddSpanTags(span, map[string]string{
        "user.id": "12345",
        "action":  "processing",
    })

    // Add events (logs)
    trace.AddSpanEvents(span, "processing-started", map[string]string{
        "timestamp": "2024-01-01T00:00:00Z",
    })

    // Simulate work...
    err := doWork(ctx)
    if err != nil {
        trace.AddSpanError(span, err)
        trace.FailSpan(span, "processing failed")
        return err
    }

    return nil
}

func doWork(ctx context.Context) error {
    // Use existing span from context instead of creating nested spans
    span := trace.SpanFromContext(ctx)

    trace.AddSpanTags(span, map[string]string{
        "step": "work-execution",
    })

    // Your business logic here
    return nil
}
```

### Production Setup with OTLP

```go
package main

import (
    "log"

    "go.pixelfactory.io/pkg/observability/trace"
)

func main() {
    // Configure via options
    provider, err := trace.NewProvider(
        trace.WithTraceEnabled(true),
        trace.WithServiceName("production-service"),
        trace.WithServiceVersion("v2.1.0"),
        trace.WithSpanExporterEndpoint("collector.example.com:4317"),
        trace.WithSpanExporterInsecure(false),
        trace.WithPropagators([]string{"tracecontext", "baggage"}),
        trace.WithHeaders(map[string]string{
            "api-key": "your-api-key",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer provider.Shutdown()

    // Application code
}
```

## Configuration

All configuration options can be set via environment variables or programmatically:

| Environment Variable | Option Function | Default | Description |
|---------------------|-----------------|---------|-------------|
| `OTEL_TRACE_ENABLED` | `WithTraceEnabled()` | `false` | Enable/disable tracing |
| `OTEL_SERVICE_NAME` | `WithServiceName()` | - | Service name for traces |
| `OTEL_SERVICE_VERSION` | `WithServiceVersion()` | `unknown` | Service version |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` | `WithSpanExporterEndpoint()` | `http://localhost:4317` | OTLP collector endpoint |
| `OTEL_EXPORTER_OTLP_TRACES_INSECURE` | `WithSpanExporterInsecure()` | `false` | Use insecure connection |
| `OTEL_EXPORTER_OTLP_HEADERS` | `WithHeaders()` | - | Custom headers for OTLP |
| `OTEL_PROPAGATORS` | `WithPropagators()` | `b3` | Propagator types (b3, tracecontext, baggage, ottrace) |

### Example with Environment Variables

```bash
export OTEL_TRACE_ENABLED=true
export OTEL_SERVICE_NAME=my-app
export OTEL_SERVICE_VERSION=1.2.3
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=collector.example.com:4317
export OTEL_PROPAGATORS=tracecontext,baggage

# Run your application
./my-app
```

## Development

### Prerequisites

- Go 1.20 or higher
- Make (optional, for convenience targets)

### Building

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build
make build
```

### Running Examples

```bash
# HTTP server example
go run example/http/server/main.go

# HTTP client example (in another terminal)
go run example/http/client/main.go
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go)
- [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector)

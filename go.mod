module go.pixelfactory.io/pkg/observability/trace

go 1.17

require (
	github.com/sethvargo/go-envconfig v0.5.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.30.0
	go.opentelemetry.io/contrib/propagators/b3 v1.5.0
	go.opentelemetry.io/contrib/propagators/ot v1.5.0
	go.opentelemetry.io/otel v1.6.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.6.0
	go.opentelemetry.io/otel/sdk v1.6.0
	go.opentelemetry.io/otel/trace v1.6.0
	google.golang.org/grpc v1.45.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.6.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/metric v0.27.0 // indirect
	go.opentelemetry.io/proto/otlp v0.12.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

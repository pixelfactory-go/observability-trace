package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.pixelfactory.io/pkg/observability/trace"
)

type helloHandler struct{}

func sayHello(ctx context.Context, name string) string {
	span := trace.SpanFromContext(ctx)
	span.SetName("sayHello")
	defer span.End()

	return fmt.Sprintf("Hello %s", name)
}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create the parent span.
	ctx, span := trace.NewSpan(ctx, "helloHandler.ServeHTTP", nil)
	defer span.End()

	resp := sayHello(ctx, "Amine")
	w.Write([]byte(resp))
}

func main() {
	// console exporter.
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	// trace provider
	prv, err := trace.NewProvider(
		trace.WithTraceEnabled(true),
		trace.WithTraceExporter(exp),
		trace.WithServiceName("server"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer prv.Shutdown()

	h := helloHandler{}
	oh := trace.HTTPHandler(h, "helloHandler")

	mux := http.NewServeMux()
	mux.Handle("/hello", oh)

	err = http.ListenAndServe(":7777", mux)
	if err != nil {
		panic(err)
	}
}

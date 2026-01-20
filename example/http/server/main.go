package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	if _, err := w.Write([]byte(resp)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
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
	defer func() {
		if shutdownErr := prv.Shutdown(); shutdownErr != nil {
			log.Printf("Failed to shutdown provider: %v", shutdownErr)
		}
	}()

	h := helloHandler{}
	oh := trace.HTTPHandler(h, "helloHandler")

	mux := http.NewServeMux()
	mux.Handle("/hello", oh)

	const (
		readTimeout  = 5 * time.Second
		writeTimeout = 10 * time.Second
		idleTimeout  = 15 * time.Second
	)
	srv := &http.Server{
		Addr:         ":7777",
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

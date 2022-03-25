package main

import (
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.pixelfactory.io/pkg/observability/trace"
)

func main() {
	// console exporter.
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	// trace provider
	prv, err := trace.NewProvider(
		trace.WithTraceExporter(exp),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer prv.Shutdown()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!\n")
	}

	// wrap http handler
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")

	http.Handle("/hello", otelHandler)
	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}

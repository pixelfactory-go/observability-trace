package main

import (
	"log"
	"net/http"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.pixelfactory.io/pkg/observability/trace"
)

type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
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
	otelth := trace.HTTPHandler(h, "helloHandler")

	mux := http.NewServeMux()
	mux.Handle("/hello", otelth)

	err = http.ListenAndServe(":7777", mux)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.pixelfactory.io/pkg/observability/trace"
)

// NewHTTPClient creates configurable http client
func NewHTTPClient(timeout time.Duration, transport *http.Transport) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

func main() {
	// console exporter.
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	// trace provider.
	prv, err := trace.NewProvider(
		trace.WithTraceEnabled(true),
		trace.WithTraceExporter(exp),
		trace.WithServiceName("client"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer prv.Shutdown()

	client := &http.Client{
		Timeout:   time.Duration(1) * time.Second,
		Transport: trace.HTTPClientTransporter(http.DefaultTransport),
	}

	resp, err := client.Get("http://localhost:7777/hello")
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body : %s", body)
}

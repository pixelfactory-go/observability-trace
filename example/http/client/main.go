package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.pixelfactory.io/pkg/observability/trace"
)

// NewHTTPClient creates configurable http client.
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
	defer func() {
		if shutdownErr := prv.Shutdown(); shutdownErr != nil {
			log.Printf("Failed to shutdown provider: %v", shutdownErr)
		}
	}()

	client := &http.Client{
		Timeout:   time.Duration(1) * time.Second,
		Transport: trace.HTTPClientTransporter(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:7777/hello", nil)
	if err != nil {
		fmt.Printf("Error %s", err) //nolint:forbidigo // Example code
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error %s", err) //nolint:forbidigo // Example code
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Body : %s", body) //nolint:forbidigo // Example code
}

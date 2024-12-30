package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	counter int
	mu      sync.Mutex
	tracer  = otel.Tracer("example-tracer")
)

// CounterHandler handles the counter API
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "CounterHandler")
	defer span.End()

	mu.Lock()
	counter++
	mu.Unlock()

	response := map[string]int{"counter": counter}
	json.NewEncoder(w).Encode(response)

	// Add an event with wrapped attributes
	span.AddEvent("Counter incremented", trace.WithAttributes(attribute.Int("counter_value", counter)))
}

// TimeHandler handles the time API
func TimeHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "TimeHandler")
	defer span.End()

	response := map[string]string{"time": time.Now().UTC().Format(time.RFC3339)}
	json.NewEncoder(w).Encode(response)

	// Add an event with wrapped attributes
	span.AddEvent("Time retrieved", trace.WithAttributes(attribute.String("timestamp", time.Now().String())))
}

// initTracer initializes the OpenTelemetry tracer
func initTracer() func() {
	ctx := context.Background()

	// Create an OTLP exporter to send traces to OTEL Collector
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector:4318"))
	if err != nil {
		fmt.Println("Failed to create OTLP exporter:", err)
		os.Exit(1)
	}

	// Create a resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", "counter-app"),
		),
	)
	if err != nil {
		fmt.Println("Failed to create resource:", err)
		os.Exit(1)
	}

	// Create a trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Println("Error shutting down tracer provider:", err)
		}
	}
}

func main() {
	shutdown := initTracer()
	if shutdown != nil {
		defer shutdown()
	}

	// Wrap handlers with OpenTelemetry instrumentation
	http.Handle("/v1/counter", otelhttp.NewHandler(http.HandlerFunc(CounterHandler), "CounterHandler"))
	http.Handle("/v1/time", otelhttp.NewHandler(http.HandlerFunc(TimeHandler), "TimeHandler"))

	fmt.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

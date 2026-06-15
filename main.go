package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// initTracer initializes the OpenTelemetry SDK and configures the OTLP gRPC exporter
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Configure the exporter to point to the local OpenTelemetry Collector service in K8s
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("opentelemetry-collector.default.svc.cluster.local:4317"),
	)
	if err != nil {
		return nil, err
	}

	// Build the Tracer Provider with a batch processing architecture
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Capture 100% of traces for the lab
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-go-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	// Initialize the pipeline
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error terminating TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("my-tracer")

	fmt.Println("Starting continuous trace generator loop...")

	// Infinite loop generating structured telemetry data every 2 seconds
	for {
		// Create the parent root trace span
		loopCtx, parentSpan := tracer.Start(context.Background(), "transaction-root-process")
		fmt.Println("\n--- New Transaction Started ---")
		fmt.Println("Processing root execution logic...")
		time.Sleep(100 * time.Millisecond) 

		// Call sub-operations passing the 'loopCtx' to link spans together
		simulateDatabaseQuery(loopCtx, tracer)
		simulateExternalAPICall(loopCtx, tracer)

		// Finish the parent span execution lifecycle
		parentSpan.End()
		fmt.Println("Trace hierarchy successfully dispatched to OTel Collector.")

		// Cool down window before triggering the next continuous sequence
		time.Sleep(2 * time.Second)
	}
}

// simulateDatabaseQuery represents an internal relational call linked to the parent context
func simulateDatabaseQuery(ctx context.Context, tracer trace.Tracer) {
	_, childSpan := tracer.Start(ctx, "SELECT * FROM orders_cache")
	defer childSpan.End()

	fmt.Println("Executing database dependency query...")
	time.Sleep(250 * time.Millisecond) // Simulating query latency
}

// simulateExternalAPICall represents a network call dependency linked to the parent context
func simulateExternalAPICall(ctx context.Context, tracer trace.Tracer) {
	_, childSpan := tracer.Start(ctx, "HTTP POST /v1/payments/charge")
	defer childSpan.End()

	fmt.Println("Calling external payment gateway API...")
	time.Sleep(150 * time.Millisecond) // Simulating round-trip latency
}

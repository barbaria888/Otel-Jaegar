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
)

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Configure the exporter to point to the local OTel Collector
	// Inside K8s, this points to your OTel collector service
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("opentelemetry-collector.default.svc.cluster.local:4317"),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
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

	// Fixed: Replaced 'ctx' with '_' to avoid the "declared and not used" compile failure
	_, span := tracer.Start(context.Background(), "simulation-span")
	defer span.End()

	fmt.Println("Application is doing heavy work inside Kubernetes...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Work done! Trace sent.")
}

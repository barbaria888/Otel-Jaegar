#!/bin/bash
set -e

rm -f go.mod go.sum

go mod init my-go-app


go get go.opentelemetry.io/otel@v1.24.0
go get go.opentelemetry.io/otel/sdk@v1.24.0
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@v1.24.0


go mod tidy


kubectl apply -f jaegar-deployment.yaml


kubectl apply -f otel-collector.yaml

echo "☸️ Step 6: Deploying Live Looping Go Trace Generator Application..."
kubectl apply -f go-app-deployment.yaml



cd ~/Otel-Jaegar
go mod init my-go-app

go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc

go mod tidy

docker build -t my-go-app:v1 .

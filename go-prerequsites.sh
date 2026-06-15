
rm -f go.mod go.sum

go mod init my-go-app


go get go.opentelemetry.io/otel@v1.24.0
go get go.opentelemetry.io/otel/sdk@v1.24.0
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@v1.24.0


go mod tidy


docker build -t my-go-app:v1 .
docker run my-go-app

module github.com/cyverse-de/email-requests

go 1.16

require (
	github.com/DavidGamba/go-getoptions v0.20.2
	github.com/cyverse-de/configurate v0.0.0-20200527185205-4e1e92866cee
	github.com/cyverse-de/messaging/v9 v9.1.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/streadway/amqp v1.0.1-0.20200716223359-e6b33f460591
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.31.0
	go.opentelemetry.io/otel v1.6.1
	go.opentelemetry.io/otel/exporters/jaeger v1.6.1
	go.opentelemetry.io/otel/sdk v1.6.1
)

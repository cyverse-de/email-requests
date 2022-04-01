package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DavidGamba/go-getoptions"
	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/go-mod/otelutils"
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const serviceName = "email-requests"

var log = logrus.WithFields(logrus.Fields{"service": serviceName})
var httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

// commandLineOptionValues represents the values of the options that were passed on the command line when this
// service was invoked.
type commandLineOptionValues struct {
	Config string
}

// parseCommandLine parses the command line and returns an options structure containing command-line options and
// parameters.
func parseCommandLine() *commandLineOptionValues {
	optionValues := &commandLineOptionValues{}
	opt := getoptions.New()

	// Default option values.
	defaultConfigPath := "/etc/iplant/de/jobservices.yml"

	// Define the command-line options.
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.StringVar(&optionValues.Config, "config", defaultConfigPath,
		opt.Alias("c"),
		opt.Description("the path to the configuration file"))

	// Parse the command line, handling requests for help and usage errors.
	_, err := opt.Parse(os.Args[1:])
	if opt.Called("help") {
		fmt.Fprint(os.Stderr, opt.Help())
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		fmt.Fprint(os.Stderr, opt.Help(getoptions.HelpSynopsis))
		os.Exit(1)
	}

	return optionValues
}

// AMQPSettings represents the settings that we require in order to connect to the AMQP exchange.
type AMQPSettings struct {
	URI          string
	ExchangeName string
	ExchangeType string
}

func main() {
	var tracerCtx, cancel = context.WithCancel(context.Background())
	defer cancel()
	shutdown := otelutils.TracerProviderFromEnv(tracerCtx, serviceName, func(e error) { log.Fatal(e) })
	defer shutdown()

	// Parse the command line.
	optionValues := parseCommandLine()

	// Load the configuration.
	cfg, err := configurate.InitDefaults(optionValues.Config, configurate.JobServicesDefaults)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the AMQP settings.
	amqpSettings := &AMQPSettings{
		URI:          cfg.GetString("amqp.uri"),
		ExchangeName: cfg.GetString("amqp.exchange.name"),
		ExchangeType: cfg.GetString("amqp.exchange.type"),
	}

	// Retrieve the base URL for the cyverse-email service.
	cyverseEmailBaseURL := cfg.GetString("iplant_email.base")

	// Create the message handler.
	handler := NewHandler(cyverseEmailBaseURL)

	// Create the message listener.
	listener, err := NewListener(handler, amqpSettings)
	if err != nil {
		log.Fatal(err)
	}

	// Listen for incoming messages.
	listener.Listen()

	// Spin until someone kills the process.
	spinner := make(chan int)
	<-spinner
}

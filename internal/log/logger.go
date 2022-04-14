package log

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

// LoggerService Handles logging information to the screen and passing to prometheus metrics when needed
type LoggerService interface {
	// Log Write message and log level to screen if maximum log level isn't exceeded
	Log(level logLevel, message string)

	// LogLoginRequest Sends metrics to prometheus for login request
	LogLoginRequest()

	// LogRegisterRequest Sends metrics to prometheus for register request
	LogRegisterRequest()
}

type logger struct {
	timeFormat      string
	maxLogLevel     int
	loginRequest    prometheus.Counter
	registerRequest prometheus.Counter
}

type logLevel struct {
	Name  string
	Level int
}

var (
	Error   = logLevel{"Error", 0}
	Warning = logLevel{"Warning", 1}
	Info    = logLevel{"Info", 2}
	Debug   = logLevel{"Debug", 3}
	Verbose = logLevel{"Verbose", 4}
)

// NewLogger Creates a new logger
func NewLogger(max logLevel, format string) LoggerService {
	if format == "" {
		format = "2012-11-01 15:04:05"
	}

	return logger{
		timeFormat:  format,
		maxLogLevel: max.Level,
		loginRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "successful_login_requests_total",
			Help:      "The total number of successful login requests",
			Namespace: "ShatteredRealms",
		}),
		registerRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "successful_register_requests_total",
			Help:      "The total number of successful registration requests",
			Namespace: "ShatteredRealms",
		}),
	}
}

// Log Writes the message to standard out with the current time, log level and message if the logger max log level is
// greater than or equal to the given log level
func (l logger) Log(level logLevel, message string) {
	if level.Level <= l.maxLogLevel {
		fmt.Printf("%s | %s: %s\n", l.formattedTime(), level.Name, message)
	}
}

// LogLoginRequest Increments the prometheus metric counter for successful login requests
func (l logger) LogLoginRequest() {
	l.loginRequest.Inc()
}

// LogRegisterRequest Increments the prometheus metric counter for successful registration requests
func (l logger) LogRegisterRequest() {
	l.registerRequest.Inc()
}

// Formats the current time to the logger time format
func (l logger) formattedTime() string {
	return time.Now().Format(l.timeFormat)
}

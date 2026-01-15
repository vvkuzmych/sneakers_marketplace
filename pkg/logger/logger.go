package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger
type Logger struct {
	logger zerolog.Logger
}

// Config holds logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	TimeFormat string
	Output     io.Writer
}

// New creates a new logger
func New(cfg Config) *Logger {
	// Set default output
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	// Parse log level
	level := parseLevel(cfg.Level)

	// Create zerolog logger
	var zlog zerolog.Logger

	if cfg.Format == "console" {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{
			Out:        cfg.Output,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		zlog = zerolog.New(output).Level(level).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		zerolog.TimeFieldFormat = time.RFC3339
		zlog = zerolog.New(cfg.Output).Level(level).With().Timestamp().Caller().Logger()
	}

	return &Logger{logger: zlog}
}

// NewDefault creates a logger with default settings
func NewDefault() *Logger {
	return New(Config{
		Level:  "info",
		Format: "console",
		Output: os.Stdout,
	})
}

// NewDevelopment creates a logger for development
func NewDevelopment() *Logger {
	return New(Config{
		Level:  "debug",
		Format: "console",
		Output: os.Stdout,
	})
}

// NewProduction creates a logger for production
func NewProduction() *Logger {
	return New(Config{
		Level:  "info",
		Format: "json",
		Output: os.Stdout,
	})
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Msgf(format, v...)
}

// With adds fields to logger
func (l *Logger) With() *Event {
	return &Event{event: l.logger.With()}
}

// WithError adds error field to logger
func (l *Logger) WithError(err error) *Event {
	return &Event{event: l.logger.With().Err(err)}
}

// WithField adds a single field to logger
func (l *Logger) WithField(key string, value interface{}) *Event {
	return &Event{event: l.logger.With().Interface(key, value)}
}

// WithFields adds multiple fields to logger
func (l *Logger) WithFields(fields map[string]interface{}) *Event {
	ctx := l.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Event{event: ctx}
}

// Event wraps zerolog.Context for chaining
type Event struct {
	event zerolog.Context
}

// Str adds string field
func (e *Event) Str(key, val string) *Event {
	e.event = e.event.Str(key, val)
	return e
}

// Int adds int field
func (e *Event) Int(key string, val int) *Event {
	e.event = e.event.Int(key, val)
	return e
}

// Int64 adds int64 field
func (e *Event) Int64(key string, val int64) *Event {
	e.event = e.event.Int64(key, val)
	return e
}

// Bool adds bool field
func (e *Event) Bool(key string, val bool) *Event {
	e.event = e.event.Bool(key, val)
	return e
}

// Err adds error field
func (e *Event) Err(err error) *Event {
	e.event = e.event.Err(err)
	return e
}

// Logger returns logger with added fields
func (e *Event) Logger() *Logger {
	return &Logger{logger: e.event.Logger()}
}

// Helper functions

func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// Global logger instance
var globalLogger *Logger

func init() {
	globalLogger = NewDefault()
}

// SetGlobal sets the global logger
func SetGlobal(l *Logger) {
	globalLogger = l
}

// Global logging functions (use global logger)

func Debug(msg string) {
	globalLogger.Debug(msg)
}

func Debugf(format string, v ...interface{}) {
	globalLogger.Debugf(format, v...)
}

func Info(msg string) {
	globalLogger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	globalLogger.Infof(format, v...)
}

func Warn(msg string) {
	globalLogger.Warn(msg)
}

func Warnf(format string, v ...interface{}) {
	globalLogger.Warnf(format, v...)
}

func Error(msg string) {
	globalLogger.Error(msg)
}

func Errorf(format string, v ...interface{}) {
	globalLogger.Errorf(format, v...)
}

func Fatal(msg string) {
	globalLogger.Fatal(msg)
}

func Fatalf(format string, v ...interface{}) {
	globalLogger.Fatalf(format, v...)
}

func WithError(err error) *Event {
	return globalLogger.WithError(err)
}

func WithField(key string, value interface{}) *Event {
	return globalLogger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) *Event {
	return globalLogger.WithFields(fields)
}

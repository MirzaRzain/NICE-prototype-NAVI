package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

type loggerOption struct {
	log zerolog.Logger
}

//nolint:gochecknoglobals // shared logger instance used across the app
var defaultLogger = &loggerOption{
	log: zerolog.New(os.Stdout),
}

var logFile *os.File

func NewLogger() {
	logDir := getenv("LOG_DIR", "logs")
	logName := getenv("LOG_FILE", "app.txt")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory %s: %v\n", logDir, err)
	}

	var output io.Writer = os.Stdout
	file, err := os.OpenFile(filepath.Join(logDir, logName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
	} else {
		logFile = file
		output = io.MultiWriter(os.Stdout, file)
	}

	log := zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).
		Logger()

	defaultLogger = &loggerOption{
		log: log,
	}
}

// Close releases the log file opened by NewLogger. It is safe to call more than once.
func Close() error {
	if logFile == nil {
		return nil
	}

	err := logFile.Close()
	logFile = nil
	return err
}

func WithContext(ctx context.Context) *Logger {
	log := defaultLogger.log

	if requestID := requestIDFromContext(ctx); requestID != "" {
		log = log.With().Str("request_id", requestID).Logger()
	}

	return &Logger{log: log}
}

func (l *Logger) Debug(args ...interface{}) *zerolog.Event {
	return l.logMessage(l.log.Debug(), args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debug().Msgf(format, args...)
}

func (l *Logger) Info(args ...interface{}) *zerolog.Event {
	return l.logMessage(l.log.Info(), args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Info().Msgf(format, args...)
}

func (l *Logger) Warn(args ...interface{}) *zerolog.Event {
	return l.logMessage(l.log.Warn(), args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warn().Msgf(format, args...)
}

func (l *Logger) Error(args ...interface{}) *zerolog.Event {
	event := l.log.Error()

	if len(args) == 0 {
		return event
	}

	if firstErr, ok := args[0].(error); ok {
		event = event.Err(firstErr)
		args = args[1:]
	}

	if len(args) == 0 {
		return event
	}

	msg, fields := splitMessageFields(args...)
	l.appendFields(event, fields)
	event.Msg(msg)

	return event
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Error().Msgf(format, args...)
}

func (l *Logger) NewError(err, newErr error, fields ...interface{}) error {
	event := l.log.Error().Err(err)
	l.appendFields(event, fields)
	event.Msg(newErr.Error())

	return newErr
}

func (l *Logger) Fatal(args ...interface{}) *zerolog.Event {
	return l.logMessage(l.log.Fatal(), args...)
}

func (l *Logger) Raw() *zerolog.Logger {
	return &l.log
}

func Debug(args ...interface{}) {
	WithContext(context.Background()).Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	WithContext(context.Background()).Debugf(format, args...)
}

func Info(args ...interface{}) {
	WithContext(context.Background()).Info(args...)
}

func Infof(format string, args ...interface{}) {
	WithContext(context.Background()).Infof(format, args...)
}

func Warn(args ...interface{}) {
	WithContext(context.Background()).Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	WithContext(context.Background()).Warnf(format, args...)
}

func Error(args ...interface{}) {
	WithContext(context.Background()).Error(args...)
}

func Errorf(format string, args ...interface{}) {
	WithContext(context.Background()).Errorf(format, args...)
}

func Errorw(msg string, err error, fields ...interface{}) {
	args := make([]interface{}, 0, len(fields)+2)
	args = append(args, err, msg)
	args = append(args, fields...)
	WithContext(context.Background()).Error(args...)
}

func PrintDebug(args ...interface{}) {
	Debug(args...)
}

func PrintInfo(args ...interface{}) {
	Info(args...)
}

func PrintWarn(args ...interface{}) {
	Warn(args...)
}

func PrintError(args ...interface{}) {
	Error(args...)
}

func PrintNewError(err, newErr error, fields ...interface{}) error {
	return WithContext(context.Background()).NewError(err, newErr, fields...)
}

func PrintFatal(err, newErr error, fields ...interface{}) {
	args := make([]interface{}, 0, len(fields)+2)
	args = append(args, err, newErr.Error())
	args = append(args, fields...)
	WithContext(context.Background()).Fatal(args...)
}

func (l *Logger) appendFields(event *zerolog.Event, fields []interface{}) {
	for i := 0; i+1 < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}

		event.Interface(key, fields[i+1])
	}
}

func requestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if value := ctx.Value("request_id"); value != nil {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}

	return ""
}

func (l *Logger) logMessage(event *zerolog.Event, args ...interface{}) *zerolog.Event {
	if len(args) == 0 {
		return event
	}

	if firstErr, ok := args[0].(error); ok {
		event = event.Err(firstErr)
		args = args[1:]
		if len(args) == 0 {
			return event
		}
	}

	msg, fields := splitMessageFields(args...)
	l.appendFields(event, fields)
	event.Msg(msg)

	return event
}

func splitMessageFields(args ...interface{}) (string, []interface{}) {
	if len(args) == 0 {
		return "", nil
	}

	msg, ok := args[0].(string)
	if !ok {
		msg = fmt.Sprint(args[0])
	}

	if len(args) == 1 {
		return msg, nil
	}

	return msg, args[1:]
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

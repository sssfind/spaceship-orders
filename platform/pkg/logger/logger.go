package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// Config Содержит параметры инициализации логгера
type Config struct {
	Level                 string
	AsJSON                bool
	ServiceName           string
	Outputs               []string // например ["stdout", "otlp"]
	OtelCollectorEndpoint string   // например "otel-collector:4317"
}

// Глобальный singleton логгер
var (
	globalLogger   *logger
	initOnce       sync.Once
	dynamicLevel   zap.AtomicLevel
	loggerProvider *sdklog.LoggerProvider
)

// logger обёртка над zap.Logger с enrich поддержкой контекста
type logger struct {
	zapLogger *zap.Logger
}

// Init сохраняет обратную совместимость для базовой инициализации
func Init(levelStr string, asJSON bool) error {
	return InitWithConfig(Config{
		Level:   levelStr,
		AsJSON:  asJSON,
		Outputs: []string{"stdout"},
	})
}

// InitWithConfig инициализирует логгер с поддержкой нескольких выходов (stdout + OTLP)
func InitWithConfig(cfg Config) error {
	var initErr error
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(cfg.Level))

		var cores []zapcore.Core

		outputsMap := make(map[string]bool)
		for _, out := range cfg.Outputs {
			outputsMap[strings.ToLower(strings.TrimSpace(out))] = true
		}
		if len(outputsMap) == 0 {
			outputsMap["stdout"] = true
		}

		// Core 1: JSON / Console -> stdout
		if outputsMap["stdout"] {
			encoderCfg := buildProductionEncoderConfig()
			var encoder zapcore.Encoder
			if cfg.AsJSON {
				encoder = zapcore.NewJSONEncoder(encoderCfg)
			} else {
				encoder = zapcore.NewConsoleEncoder(encoderCfg)
			}

			stdoutCore := zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stdout),
				dynamicLevel,
			)
			cores = append(cores, stdoutCore)
		}

		// Core 2: OTLP Log Exporter -> OTEL Collector (OTLP/gRPC)
		if outputsMap["otlp"] && cfg.OtelCollectorEndpoint != "" {
			ctx := context.Background()

			res, err := resource.New(ctx,
				resource.WithAttributes(
					semconv.ServiceNameKey.String(cfg.ServiceName),
				),
			)
			if err != nil {
				initErr = err
				return
			}

			exporter, err := otlploggrpc.New(ctx,
				otlploggrpc.WithEndpoint(cfg.OtelCollectorEndpoint),
				otlploggrpc.WithInsecure(),
			)
			if err != nil {
				initErr = err
				return
			}

			loggerProvider = sdklog.NewLoggerProvider(
				sdklog.WithResource(res),
				sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
			)

			otlpCore := otelzap.NewCore(
				cfg.ServiceName,
				otelzap.WithLoggerProvider(loggerProvider),
			)
			cores = append(cores, otlpCore)
		}

		combinedCore := zapcore.NewTee(cores...)
		zapLogger := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	return initErr
}

func Shutdown(ctx context.Context) error {
	if globalLogger != nil && globalLogger.zapLogger != nil {
		_ = globalLogger.zapLogger.Sync()
	}
	if loggerProvider != nil {
		return loggerProvider.Shutdown(ctx)
	}
	return nil
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // время
		LevelKey:       "level",                     // уровень логирования
		NameKey:        "logger",                    // имя логгера, если используется
		CallerKey:      "caller",                    // откуда вызван лог
		MessageKey:     "message",                   // текст сообщения
		StacktraceKey:  "stacktrace",                // стектрейс для ошибок
		LineEnding:     zapcore.DefaultLineEnding,   // перенос строки
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // читаемый ISO 8601 формат
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // короткий caller
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// SetLevel динамически меняет уровень логирования
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

func InitForBenchmark() {
	core := zapcore.NewNopCore()

	globalLogger = &logger{
		zapLogger: zap.New(core),
	}
}

// logger возвращает глобальный enrich-aware логгер
func Logger() *logger {
	return globalLogger
}

// NopLogger устанавливает глобальный логгер в no-op режим
func SetNopLogger() {
	globalLogger = &logger{
		zapLogger: zap.NewNop(),
	}
}

// Sync сбрасывает буферы логгера
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext создает enrich-aware логгер с контекстом
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug enrich-aware debug log
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

// Info enrich-aware info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

// Warn enrich-aware warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

// Error enrich-aware error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

// Fatal enrich-aware fatal log
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

// Instance methods для enrich loggers

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// parseLevel конвертирует строковый уровень в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext вытаскивает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

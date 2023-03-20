package api

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZap(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapCfg := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	return zapCfg.Build()
}

type MuxLoggingMiddleware struct {
	logger *zap.Logger
}

func NewMuxLoggingMiddleware(logger *zap.Logger) *MuxLoggingMiddleware {
	return &MuxLoggingMiddleware{logger}
}

func (m *MuxLoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		statusCode := "0"

		wrw, ok := w.(*wrappedResponseWriter)
		if ok {
			statusCode = fmt.Sprintf("%d", wrw.statusCode)
		}

		m.logger.Debug(
			"request",
			zap.String("proto", r.Proto),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.String("statusCode", statusCode),
		)
	})
}

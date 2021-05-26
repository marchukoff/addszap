package addszap

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(debug bool) *zap.Logger {
	return newLogger(debug, os.Stdout)
}

func NewWriteLogger(debug bool, writers ...io.Writer) *zap.Logger {
	var writer io.Writer

	switch len(writers) {
	case 0:
		writer = os.Stdout
	case 1:
		writer = writers[0]
	default:
		syncers := make([]io.Writer, len(writers))
		for i, w := range writers {
			syncers[i] = zapcore.AddSync(w)
		}
		writer = io.MultiWriter(syncers...)
	}

	return newLogger(debug, zapcore.AddSync(writer))
}

func newLogger(debug bool, w io.Writer) *zap.Logger {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "time"
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w),
		level,
	)

	return zap.New(core)
}

package logger

import (
	"context"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Infof(ctx context.Context, format string, msg ...interface{})
	Errorf(ctx context.Context, format string, msg ...interface{})
	Warnf(ctx context.Context, format string, msg ...interface{})
	Debugf(ctx context.Context, format string, msg ...interface{})
}

type appLogger struct {
	internal *logrus.Logger
}

var (
	Log  Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		l := logrus.New()

		l.Out = os.Stdout

		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})

		l.SetLevel(logrus.InfoLevel)

		Log = &appLogger{
			internal: l,
		}
	})
}

func (l *appLogger) build(ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(l.internal)

	if ctx != nil {
		if reqID, ok := ctx.Value("request_id").(string); ok {
			entry = entry.WithField("request_id", reqID)
		}
	}

	return entry
}

func (l *appLogger) Infof(ctx context.Context, format string, msg ...interface{}) {
	l.build(ctx).Infof(format, msg...)
}

func (l *appLogger) Errorf(ctx context.Context, format string, msg ...interface{}) {
	l.build(ctx).Errorf(format, msg...)
}

func (l *appLogger) Warnf(ctx context.Context, format string, msg ...interface{}) {
	l.build(ctx).Warnf(format, msg...)
}

func (l *appLogger) Debugf(ctx context.Context, format string, msg ...interface{}) {
	l.build(ctx).Debugf(format, msg...)
}

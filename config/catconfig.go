package config

import (
	"fmt"

	"go.uber.org/zap"
	configcat "gopkg.in/configcat/go-sdk.v1"
)

type catLogger struct {
	l *zap.Logger
}

func (l *catLogger) Prefix(p string) configcat.Logger {
	return &catLogger{l.l.With(zap.String("scope", p))}
}

func (l *catLogger) Print(v ...interface{})            { l.l.Debug(fmt.Sprint(v...)) }
func (l *catLogger) Printf(f string, v ...interface{}) { l.l.Debug(fmt.Sprintf(f, v...)) }

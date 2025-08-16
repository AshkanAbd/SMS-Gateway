package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var inner = log.Logger
var instance = &zeroLog{
	inner,
}

func SetLogLevel(lvlStr string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var err error
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		lvl = zerolog.DebugLevel
	}

	inner.Level(lvl)
}

type Logger interface {
	Trace(msg string)
	Tracef(format string, data ...any)
	Debug(msg string)
	Debugf(format string, data ...any)
	Info(msg string)
	Infof(format string, data ...any)
	Warn(msg string)
	Warnf(format string, data ...any)
	Error(err error, msg string)
	Errorf(err error, format string, data ...any)
	Fatal(err error, msg string)
	Fatalf(err error, format string, data ...any)
}

func Trace(format string, data ...any) {
	if len(data) == 0 {
		instance.Trace(format)
		return
	}
	instance.Tracef(format, data...)
}

func Debug(format string, data ...any) {
	if len(data) == 0 {
		instance.Debug(format)
		return
	}
	instance.Debugf(format, data...)
}

func Info(format string, data ...any) {
	if len(data) == 0 {
		instance.Info(format)
		return
	}
	instance.Infof(format, data...)
}

func Warn(format string, data ...any) {
	if len(data) == 0 {
		instance.Warn(format)
		return
	}
	instance.Warnf(format, data...)
}

func Error(err error, format string, data ...any) {
	if len(data) == 0 {
		instance.Error(err, format)
		return
	}
	instance.Errorf(err, format, data...)
}

func Fatal(err error, format string, data ...any) {
	if len(data) == 0 {
		instance.Fatal(err, format)
		return
	}
	instance.Fatalf(err, format, data...)
}

type zeroLog struct {
	logger zerolog.Logger
}

func (z *zeroLog) Tracef(format string, data ...any) {
	z.logger.Trace().Msgf(format, data...)
}

func (z *zeroLog) Trace(msg string) {
	z.logger.Trace().Msg(msg)
}

func (z *zeroLog) Debugf(format string, data ...any) {
	z.logger.Debug().Msgf(format, data...)
}

func (z *zeroLog) Debug(msg string) {
	z.logger.Debug().Msg(msg)
}

func (z *zeroLog) Infof(format string, data ...any) {
	z.logger.Info().Msgf(format, data...)
}

func (z *zeroLog) Info(msg string) {
	z.logger.Info().Msg(msg)
}

func (z *zeroLog) Warnf(format string, data ...any) {
	z.logger.Warn().Msgf(format, data...)
}

func (z *zeroLog) Warn(msg string) {
	z.logger.Warn().Msg(msg)
}

func (z *zeroLog) Errorf(err error, format string, data ...any) {
	z.logger.Error().Err(err).Msgf(format, data...)
}

func (z *zeroLog) Error(err error, msg string) {
	z.logger.Error().Err(err).Msg(msg)
}

func (z *zeroLog) Fatalf(err error, format string, data ...any) {
	z.logger.Fatal().Err(err).Msgf(format, data...)
}

func (z *zeroLog) Fatal(err error, msg string) {
	z.logger.Fatal().Err(err).Msg(msg)
}

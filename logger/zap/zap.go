package zap

import (
	"context"
	"fmt"
	"github.com/qoqozhang/go-web-basic/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"sync"
)

type Zaplog struct {
	opts logger.Options
	sync.RWMutex
	zap    *zap.Logger
	cfg    *zap.Config
	fields map[string]interface{}
}

func (z *Zaplog) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&z.opts)
	}

	zapConfig := zap.Config{}
	if z.cfg != nil {
		zapConfig = *z.cfg
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevel()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	if z.opts.Level != logger.InfoLevel {
		zapConfig.Level.SetLevel(loggerToZapLevel(z.opts.Level))
	}

	// 添加多个日志输出源
	var sync []zapcore.WriteSyncer
	for _, out := range z.opts.Out {
		sync = append(sync, zapcore.AddSync(out))
	}

	logCore := zapcore.NewCore(zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
		zapcore.NewMultiWriteSyncer(sync...),
		zapConfig.Level)

	log := zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))

	if z.opts.Fields != nil {
		var data []zap.Field
		for k, v := range z.opts.Fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}
	if z.fields != nil {
		var data []zap.Field
		for k, v := range z.fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}

	z.cfg = &zapConfig
	z.zap = log
	z.fields = map[string]interface{}{}
	return nil
}

func (z *Zaplog) Options() logger.Options {
	return z.opts
}

func (z *Zaplog) Fields(fields map[string]interface{}) logger.Logger {
	z.Lock()
	nfields := make(map[string]interface{}, len(z.fields)+len(fields))
	for k, v := range z.fields {
		nfields[k] = v
	}
	z.Unlock()
	for k, v := range fields {
		nfields[k] = v
	}
	zl := &Zaplog{
		cfg:    z.cfg,
		zap:    z.zap,
		opts:   z.opts,
		fields: nfields,
	}
	return zl
}
func (z Zaplog) Log(level logger.Level, args ...interface{}) {
	z.RLock()
	data := make([]zap.Field, 0, len(z.fields))
	for k, v := range z.fields {
		data = append(data, zap.Any(k, v))
	}
	z.RUnlock()
	msg := fmt.Sprint(args...)
	lvl := loggerToZapLevel(level)
	switch lvl {
	case zap.DebugLevel:
		z.zap.Debug(msg, data...)
	case zap.InfoLevel:
		z.zap.Info(msg, data...)
	case zap.WarnLevel:
		z.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		z.zap.Error(msg, data...)
	case zap.FatalLevel:
		z.zap.Fatal(msg, data...)
	}
}
func (z Zaplog) Logf(level logger.Level, format string, args ...interface{}) {
	z.RLock()
	data := make([]zap.Field, 0, len(z.fields))
	for k, v := range z.fields {
		data = append(data, zap.Any(k, v))
	}
	z.RUnlock()
	msg := fmt.Sprintf(format, args...)
	lvl := loggerToZapLevel(level)
	switch lvl {
	case zap.DebugLevel:
		z.zap.Debug(msg, data...)
	case zap.InfoLevel:
		z.zap.Info(msg, data...)
	case zap.WarnLevel:
		z.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		z.zap.Error(msg, data...)
	case zap.FatalLevel:
		z.zap.Fatal(msg, data...)
	}

}
func (z Zaplog) String() string {
	return "zap"
}

func NewLogger(opts ...logger.Option) (logger.Logger, error) {
	options := logger.Options{
		Level:   logger.InfoLevel,
		Out:     []io.Writer{},
		Context: context.Background(),
	}
	l := &Zaplog{opts: options}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}
	return l, nil
}

func loggerToZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.DebugLevel:
		return zap.DebugLevel
	case logger.InfoLevel:
		return zap.InfoLevel
	case logger.WarnLevel:
		return zap.WarnLevel
	case logger.ErrorLevel:
		return zap.ErrorLevel
	case logger.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

package logger

import (
	"context"
	"io"
)

type Option func(*Options)

type Options struct {
	Level   Level
	Fields  map[string]interface{}
	Out     []io.Writer
	Context context.Context
	Name    string
}

func SetLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}
func SetFields(fields map[string]interface{}) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}
func SetOut(out io.Writer) Option {
	return func(o *Options) {
		o.Out = append(o.Out, out)
	}
}

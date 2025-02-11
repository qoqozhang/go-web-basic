package logger

type Logger interface {
	Init(options ...Option) error
	Options() Options
	Fields(fields map[string]interface{}) Logger
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	String() string
}

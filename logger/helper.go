package logger

type Helper struct {
	Logger
	fields map[string]interface{}
}

func NewHelper(log Logger) *Helper {
	return &Helper{Logger: log}
}

func (h *Helper) Debug(args ...interface{}) {
	basic(h, DebugLevel, args...)
}
func (h *Helper) Debugf(template string, args ...interface{}) {
	basicF(h, DebugLevel, template, args...)
}

func (h *Helper) Info(args ...interface{}) {
	basic(h, InfoLevel, args...)
}
func (h *Helper) Infof(template string, args ...interface{}) {
	basicF(h, InfoLevel, template, args...)
}
func (h *Helper) Warn(args ...interface{}) {
	basic(h, WarnLevel, args...)
}
func (h *Helper) Warnf(template string, args ...interface{}) {
	basicF(h, WarnLevel, template, args...)
}

func (h *Helper) Error(args ...interface{}) {
	basic(h, ErrorLevel, args...)
}
func (h *Helper) Errorf(template string, args ...interface{}) {
	basicF(h, ErrorLevel, template, args...)
}
func (h *Helper) Fatal(args ...interface{}) {
	basic(h, FatalLevel, args...)
}
func (h *Helper) Fatalf(template string, args ...interface{}) {
	basicF(h, FatalLevel, template, args...)
}

// WithError 向Helper的Fields字段中添加error字段
func (h *Helper) WithError(err error) *Helper {
	fields := copyFields(h.fields)
	fields["error"] = err
	return &Helper{Logger: h.Logger, fields: fields}
}

// WithFields 向Helper的Fields字段中添加字段
func (h *Helper) WithFields(fields map[string]interface{}) *Helper {
	nfields := copyFields(h.fields)
	for k, v := range fields {
		nfields[k] = v
	}
	return &Helper{Logger: h.Logger, fields: nfields}
}

// 日志输出的基础函数
func basic(h *Helper, level Level, args ...interface{}) {
	if !h.Logger.Options().Level.Enabled(level) {
		return
	}
	h.Logger.Fields(h.fields).Log(level, args...)
}

// 格式化日志输出的基础函数
func basicF(h *Helper, level Level, template string, args ...interface{}) {
	if !h.Logger.Options().Level.Enabled(level) {
		return
	}
	h.Logger.Fields(h.fields).Logf(level, template, args...)
}

// 用于拷贝helper的fields字段
func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

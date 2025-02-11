package zap

import (
	"errors"
	"github.com/qoqozhang/go-web-basic/logger"
	"github.com/qoqozhang/go-web-basic/logger/writer"
	"testing"
)

func TestNewLogger(t *testing.T) {
	zap, err := NewLogger()
	if err != nil {
		t.Errorf("create zap fail: %v\n", err)
	}
	zap.Log(logger.InfoLevel, "test log function")
	zap.Logf(logger.InfoLevel, "test logf function: value: %s", "2025-2-8")
	zap.Fields(map[string]interface{}{
		"name": "zhanghailiang",
		"age":  230,
	}).Log(logger.InfoLevel, "test log function with fields")
	zap.Fields(map[string]interface{}{
		"name": "zhanghailiang",
		"age":  22,
	}).Logf(logger.InfoLevel, "test logf function: value: %s", "2025-2-8")
	zap.Log(logger.DebugLevel, "test debug level log function")
	zap.Init(logger.SetLevel(logger.DebugLevel))
	zap.Logf(logger.DebugLevel, "test debug level log function: value: %s", "2025-2-9")
}
func TestZaploger(t *testing.T) {
	rotateFile := writer.RotateFileLog{
		FilePrefix: "user",
		FilePath:   "logs/",
		MaxAge:     90,
	}

	syslog := writer.Syslog{
		Address:  "10.1.1.236",
		Port:     514,
		Protocol: "udp",
	}
	zaplogger, _ := NewLogger(logger.SetFields(map[string]interface{}{
		"namespace": "admin",
	}),
		logger.SetOut(&rotateFile),
		logger.SetOut(&syslog),
	)
	helper := logger.NewHelper(zaplogger)
	helper.Info("test logger")
	helper.Infof("test logger with fields: %s", "zhanghailiang ")
	helper.WithFields(map[string]interface{}{
		"name": "helper",
	}).Info("test log function")
	helper.WithError(errors.New("test error")).Info("test log function with error")
	helper.WithError(errors.New("test error")).Infof("test logf function with error: %s, %d", "zzzzz", 333)

}

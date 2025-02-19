package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qoqozhang/go-web-basic/logger/writer"
)

// GinLoggerToFile  写入gin日志到日志文件中
func GinLoggerToFile(logPrefix string) gin.HandlerFunc {
	log := writer.RotateFileLog{
		FilePrefix: "server",
		MaxAge:     180,
	}
	GinLoggerConfig := gin.LoggerConfig{
		Output: &log,
	}
	return gin.LoggerWithConfig(GinLoggerConfig)
}

// GinLoggerToSyslog  写入gin日志到syslog服务中
func GinLoggerToSyslog(server string, port int, proto string) gin.HandlerFunc {
	syslog := writer.Syslog{
		Address:  server,
		Port:     port,
		Protocol: proto,
	}
	GinLoggerConfig := gin.LoggerConfig{
		Output: &syslog,
	}
	return gin.LoggerWithConfig(GinLoggerConfig)
}

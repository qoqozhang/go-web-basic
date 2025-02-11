package writer

import (
	"fmt"
	gsyslog "github.com/hashicorp/go-syslog"
	"log"
	"net"
)

// Syslog 一个实现了io.Writer接口功能的syslog服务器连接对象，不支持windows上调用
type Syslog struct {
	Address   string            `json:"address"`  // syslog服务器的IP地址
	Port      int               `json:"port"`     // syslog 端口
	Protocol  string            `json:"protocol"` //syslog协议，支持tcp、udp
	syslogger gsyslog.Syslogger `json:"-"`
}

// Write 输出日志到syslog服务器
// 如果没有连接到syslog，则初始化一个连接
func (s *Syslog) Write(b []byte) (int, error) {
	if s.syslogger == nil {
		if err := s.init(); err != nil {
			log.Fatalf("syslog init failed: %v", err)
		}
	}
	return s.syslogger.Write(b)
}

// init 初始化一个连接syslog服务器的请求
func (s *Syslog) init() error {
	ip := net.ParseIP(s.Address)
	if ip == nil {
		log.Fatalf("unvalid ip address: %s", s.Address)
	}
	syslogger, err := gsyslog.DialLogger(s.Protocol, fmt.Sprintf("%s:%d", ip.String(), s.Port), gsyslog.LOG_DEBUG, "LOCAL7", "app")
	if err != nil {
		return err
	}
	s.syslogger = syslogger
	return nil
}

func (s *Syslog) Close() error {
	return s.syslogger.Close()
}

package nsq_logger

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/youzan/go-nsq"
)

var (
	nsqDebugLevel = nsq.LogLevelDebug.String()
	nsqInfoLevel  = nsq.LogLevelInfo.String()
	nsqWarnLevel  = nsq.LogLevelWarning.String()
	nsqErrLevel   = nsq.LogLevelError.String()
)

type NSQLogrusLogger struct{}

func NewNSQLogrusLoggerAtLevel(l log.Level) (NSQLogrusLogger, nsq.LogLevel) {
	level := nsq.LogLevelWarning
	switch l {
	case log.DebugLevel:
		level = nsq.LogLevelDebug
	case log.InfoLevel:
		level = nsq.LogLevelInfo
	case log.WarnLevel:
		level = nsq.LogLevelWarning
	case log.ErrorLevel:
		level = nsq.LogLevelError
	}
	return NSQLogrusLogger{}, level
}

func (n NSQLogrusLogger) Output(_ int, s string) error {
	if len(s) > 3 {
		msg := strings.TrimSpace(s[3:])
		switch s[:3] {
		case nsqDebugLevel:
			log.Debugln(msg)
		case nsqInfoLevel:
			log.Infoln(msg)
		case nsqWarnLevel:
			log.Warnln(msg)
		case nsqErrLevel:
			log.Errorln(msg)
		default:
			log.Infoln(msg)
		}
	}
	return nil
}

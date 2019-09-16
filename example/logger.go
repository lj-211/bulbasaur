package example

import (
	"os"

	gologging "github.com/op/go-logging"
	"github.com/pkg/errors"
)

var Log Logger = Logger{}

type Logger struct {
	Logger *gologging.Logger
}

func (l Logger) Info(args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Info(args)
	}
}

func (l Logger) Infof(format string, args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Infof(format, args...)
	}
}

func (l Logger) Debug(args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Debug(args)
	}
}

func (l Logger) Debugf(format string, args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Debugf(format, args)
	}
}

func (l Logger) Warning(args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Warning(args)
	}
}

func (l Logger) Warningf(format string, args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Warningf(format, args)
	}
}

func (l Logger) Error(args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Error(args)
	}
}

func (l Logger) Errorf(format string, args ...interface{}) {
	if l.Logger != nil {
		l.Logger.Errorf(format, args)
	}
}
func (l Logger) SetLevel(lv int) {}
func (l Logger) GetLevel() int   { return 0 }

const logFormat string = "%{level}: [%{time:2006-01-02 15:04:05.000}][%{pid}][%{module}][%{shortfile}] - %{message}"

func InitGoLogging(name string, fpath string, fname string) error {
	log_fp, err := os.OpenFile(fpath+"/"+fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "打开文件失败")
	}
	err = initGoLogging(name, log_fp)
	if err != nil {
		return errors.Wrap(err, "初始化日志失败")
	}

	return nil
}

func InitGoLoggingStdout(name string) error {
	err := initGoLogging(name, os.Stdout)
	if err != nil {
		return errors.Wrap(err, "初始化日志失败")
	}

	return nil
}

func initGoLogging(name string, fd *os.File) error {
	if fd == nil {
		return errors.New("无效的日志文件句柄")
	}

	backend := gologging.NewLogBackend(fd, "", 0)
	format := gologging.MustStringFormatter(logFormat)
	formatter := gologging.NewBackendFormatter(backend, format)
	level := gologging.AddModuleLevel(formatter)
	level.SetLevel(gologging.INFO, "")

	gologging.SetBackend(level)

	Log.Logger = gologging.MustGetLogger(name)

	return nil
}

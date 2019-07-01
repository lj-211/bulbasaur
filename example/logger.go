package example

import (
	"os"

	gologging "github.com/op/go-logging"
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
		l.Logger.Infof(format, args)
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

const logFormat string = gologging.MustStringFormatter("%{level}: [%{time:2006-01-02 15:04:05.000}][%{pid}][grt_id:%{goroutineid}][grt_count:%{goroutinecount}][%{module}][%{shortfile}][%{message}]")

func initGoLogging(name string, fpath string, fname string) error {
	backend := gologging.NewLogBackend(os.Stdout, "", 0)

	log_fp, err := os.OpenFile(fpath+"/"+fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "打开文件失败")
	}

	backend := gologging.NewLogBackend(info_log_fp, "", 0)
	formatter := gologging.NewBackendFormatter(backend, format)
	level := gologging.AddModuleLevel(formatter)
	level.SetLevel(gologging.INFO, "")

	gologging.SetBackend(level)

	Log.Logger = gologging.MustGetLogger(name)

	return nil
}

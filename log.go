package bulbasaur

import (
	"github.com/lj-211/common/cmodel"
)

var Logger LoggerIns = LoggerIns{}

type LoggerIns struct {
	Log cmodel.Logger
}

func (l LoggerIns) Info(args ...interface{}) {
	if l.Log != nil {
		l.Log.Info(args)
	}
}

func (l LoggerIns) Infof(format string, args ...interface{}) {
	if l.Log != nil {
		l.Log.Infof(format, args)
	}
}

func (l LoggerIns) Debug(args ...interface{}) {
	if l.Log != nil {
		l.Log.Debug(args)
	}
}

func (l LoggerIns) Debugf(format string, args ...interface{}) {
	if l.Log != nil {
		l.Log.Debugf(format, args)
	}
}

func (l LoggerIns) Warning(args ...interface{}) {
	if l.Log != nil {
		l.Log.Warning(args)
	}
}

func (l LoggerIns) Warningf(format string, args ...interface{}) {
	if l.Log != nil {
		l.Log.Warningf(format, args)
	}
}

func (l LoggerIns) Error(args ...interface{}) {
	if l.Log != nil {
		l.Log.Error(args)
	}
}

func (l LoggerIns) Errorf(format string, args ...interface{}) {
	if l.Log != nil {
		l.Log.Errorf(format, args)
	}
}

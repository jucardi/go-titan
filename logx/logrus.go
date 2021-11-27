package logx

import "github.com/sirupsen/logrus"

type logrusWrapper struct {
	*logrus.Logger
}

func (l *logrusWrapper) WithObj(obj interface{}) IEntry {
	if obj == nil {
		return nilLogger{}
	}
	return l.fromLogrusEntry(l.Logger.WithField("obj", obj))
}

func (l *logrusWrapper) WithFields(fields map[string]interface{}) IEntry {
	return l.fromLogrusEntry(l.Logger.WithFields(fields))
}

func (l *logrusWrapper) WithField(key string, val interface{}) IEntry {
	return l.fromLogrusEntry(l.Logger.WithField(key, val))
}

func (l *logrusWrapper) SetLevel(level Level) {
	l.Logger.SetLevel(l.toLogrusLevel(level))
}

func (l *logrusWrapper) GetLevel() Level {
	return l.fromLogrusLevel(l.Logger.GetLevel())
}

func (l *logrusWrapper) toLogrusLevel(level Level) logrus.Level {
	ret := -1
	for l := level; l > 0; l = l >> 1 {
		ret++
	}
	return logrus.Level(ret)
}

func (l *logrusWrapper) fromLogrusLevel(level logrus.Level) Level {
	return Level(0x1 << level)
}

func (l *logrusWrapper) fromLogrusEntry(e *logrus.Entry, logger ...*logrus.Logger) IEntry {
	if len(logger) > 0 && logger[0] != nil {
		e.Logger = logger[0]
	}
	return &logrusEntryWrapper{
		Entry: e,
	}
}

func NewLogrus() ILogger {
	return &logrusWrapper{logrus.New()}
}

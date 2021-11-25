package logx

import "github.com/sirupsen/logrus"

type logrusEntryWrapper struct {
	*logrus.Entry
}

func (e *logrusEntryWrapper) WithObj(obj interface{}) IEntry {
	if obj == nil {
		return nilLogger{}
	}
	return e.fromLogrusEntry(e.Entry.WithField("obj", obj))
}

func (e *logrusEntryWrapper) WithFields(fields map[string]interface{}) IEntry {
	return e.fromLogrusEntry(e.Entry.WithFields(fields))
}

func (e *logrusEntryWrapper) WithField(key string, val interface{}) IEntry {
	return e.fromLogrusEntry(e.Entry.WithField(key, val))
}

func (*logrusEntryWrapper) fromLogrusEntry(e *logrus.Entry, logger ...*logrus.Logger) IEntry {
	if len(logger) > 0 && logger[0] != nil {
		e.Logger = logger[0]
	}
	return &logrusEntryWrapper{
		Entry: e,
	}
}

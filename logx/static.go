package logx

var (
	DefaultLogger = NewLogrus()
	loggers       = map[string]ILogger{}
)

func Get(name string) ILogger {
	if v, ok := loggers[name]; ok {
		return v
	}
	ret := NewLogrus()
	loggers[name] = ret
	return ret
}

func GetLevel() Level                                 { return DefaultLogger.GetLevel() }
func SetLevel(level Level)                            { DefaultLogger.SetLevel(level) }
func Trace(args ...interface{})                       { DefaultLogger.Trace(args...) }
func Tracef(format string, args ...interface{})       { DefaultLogger.Tracef(format, args...) }
func Debug(args ...interface{})                       { DefaultLogger.Debug(args...) }
func Debugf(format string, args ...interface{})       { DefaultLogger.Debugf(format, args...) }
func Info(args ...interface{})                        { DefaultLogger.Info(args...) }
func Infof(format string, args ...interface{})        { DefaultLogger.Infof(format, args...) }
func Warn(args ...interface{})                        { DefaultLogger.Warn(args...) }
func Warnf(format string, args ...interface{})        { DefaultLogger.Warnf(format, args...) }
func Error(args ...interface{})                       { DefaultLogger.Error(args...) }
func Errorf(format string, args ...interface{})       { DefaultLogger.Errorf(format, args...) }
func Fatal(args ...interface{})                       { DefaultLogger.Fatal(args...) }
func Fatalf(format string, args ...interface{})       { DefaultLogger.Fatalf(format, args...) }
func Panic(args ...interface{})                       { DefaultLogger.Panic(args...) }
func Panicf(format string, args ...interface{})       { DefaultLogger.Panicf(format, args...) }
func WithFields(fields map[string]interface{}) IEntry { return DefaultLogger.WithFields(fields) }
func WithField(key string, val interface{}) IEntry    { return DefaultLogger.WithField(key, val) }
func WithObj(obj interface{}) IEntry                  { return DefaultLogger.WithObj(obj) }

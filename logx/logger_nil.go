package logx

var (
	_ ILogger = (*nilLogger)(nil)
)

// ILogger implementation that does nothing on function calls.
type nilLogger struct {
}

func (n nilLogger) GetLevel() Level                                      { return LevelDebug }
func (n nilLogger) Name() string                                         { return "" }
func (n nilLogger) SetLevel(level Level)                                 {}
func (n nilLogger) Trace(args ...interface{})                            {}
func (n nilLogger) Debug(args ...interface{})                            {}
func (n nilLogger) Debugf(format string, args ...interface{})            {}
func (n nilLogger) Tracef(format string, args ...interface{})            {}
func (n nilLogger) Info(args ...interface{})                             {}
func (n nilLogger) Infof(format string, args ...interface{})             {}
func (n nilLogger) Warn(args ...interface{})                             {}
func (n nilLogger) Warnf(format string, args ...interface{})             {}
func (n nilLogger) Error(args ...interface{})                            {}
func (n nilLogger) Errorf(format string, args ...interface{})            {}
func (n nilLogger) Fatal(args ...interface{})                            {}
func (n nilLogger) Fatalf(format string, args ...interface{})            {}
func (n nilLogger) Panic(args ...interface{})                            {}
func (n nilLogger) Panicf(format string, args ...interface{})            {}
func (n nilLogger) Log(level Level, args ...interface{})                 {}
func (n nilLogger) Logf(level Level, format string, args ...interface{}) {}
func (n nilLogger) WithObj(obj interface{}) IEntry                       { return n }
func (n nilLogger) WithFields(fields map[string]interface{}) IEntry      { return n }
func (n nilLogger) WithField(key string, val interface{}) IEntry         { return n }

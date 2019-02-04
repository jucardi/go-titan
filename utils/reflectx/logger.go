package reflectx

type ILogger interface {
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var (
	nilLogger = &defaultLogger{}
	log       = func() ILogger { return nilLogger }
)

type defaultLogger struct{}

func (d *defaultLogger) Warnf(format string, args ...interface{})  {}
func (d *defaultLogger) Errorf(format string, args ...interface{}) {}

func SetLogger(f func() ILogger) {
	log = f
}

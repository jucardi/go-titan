package logging

import (
	"fmt"
	"net/http"

	"github.com/jucardi/go-titan/errors"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/errorx"
	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/net/rest/config"
	"github.com/jucardi/go-titan/net/rest/middleware/metrics"
)

var (
	httpDumpOn = 500
)

type dumpObj struct {
	Errors        []*errorx.Error `json:"errors,omitempty"         yaml:"errors,omitempty"`
	ErrorMessages []string        `json:"error_messages,omitempty" yaml:"error_messages,omitempty"`
	HttpRequest   string          `json:"http_request,omitempty"   yaml:"http_request,omitempty"`

	request *http.Request
}

func (d *dumpObj) Error() error {
	if len(d.Errors) > 0 && len(d.ErrorMessages) == 0 {
		return errorx.MergeErrorx(d.Errors...)
	}
	if len(d.Errors) == 0 && len(d.ErrorMessages) > 0 {
		return errors.Join("multiple errors occurred", d.ErrorMessages...)
	}
	err1 := errorx.MergeErrorx(d.Errors...)
	err2 := errors.Join("multiple errors occurred", d.ErrorMessages...)
	return errorx.MergeErrors(err1, err2)
}

func (d *dumpObj) Request() *http.Request {
	return d.request
}

func init() {
	config.AddReloadCallback(func(config *config.RestConfig) {
		if minCode := config.Reporting.MinStatus; minCode > 0 {
			httpDumpOn = minCode
		} else if minCode < 0 {
			httpDumpOn = 0
		}
	})
}

// Handler is a Gin handler which captures response times between requests and consistent logging formatting
func Handler(c *rest.Context) {
	path := c.Request.URL.Path
	c.Next()
	latency := metrics.GetMeasuredLatency(c)
	status := c.Writer.Status()

	if !shouldLog(status, path, logx.GetLevel()) {
		return
	}

	var lx = logx.WithFields(map[string]interface{}{
		"Source":  c.ClientIP(),
		"Latency": latency.String(),
	})

	msg := fmt.Sprintf(
		`[%3d] %s %s`,
		c.Writer.Status(),
		c.Request.Method,
		path)

	dump := &dumpObj{
		request: c.Request,
	}

	for _, err := range c.Errors {
		if e, ok := err.Err.(*errorx.Error); ok {
			dump.Errors = append(dump.Errors, e)
		} else {
			dump.ErrorMessages = append(dump.ErrorMessages, err.Error())
		}
	}

	if status >= httpDumpOn {
		dump.HttpRequest = c.DumpRequest()
	}

	if len(dump.ErrorMessages) > 0 || len(dump.HttpRequest) > 0 || len(dump.Errors) > 0 {
		lx = lx.WithObj(dump)
	}

	if c.Writer.Status() >= 400 && c.Writer.Status() < 500 {
		lx.Warn(msg)
	} else if c.Writer.Status() >= 500 {
		lx.Error(msg)
	} else {
		lx.Debug(msg)
	}
}

func shouldLog(status int, path string, level logx.Level) bool {
	if path == "/info" && level >= logx.LevelDebug {
		return true
	}
	switch level {
	case logx.LevelErrors:
		return status >= 500
	case logx.LevelWarn:
		return status >= 400
	}
	return true
}

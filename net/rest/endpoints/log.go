package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/errorx"
	"github.com/jucardi/go-titan/net/rest"
)

// AddLogLevel adds the `/logger/:level` endpoint to the given router.
func AddLogLevel(router *gin.Engine) {
	router.POST("/logger/:name/:level", func(context *gin.Context) {
		updateLoggerLevel(rest.NewContext(context, false))
	})
	router.POST("/loggers/:level", func(context *gin.Context) {
		updateLoggersLevel(rest.NewContext(context, false))
	})
}

// swagger:route POST /loggers/:level update loggers level
//
// Updates the log level of the service
//
// Responses:
//   200: Confirmation message
func updateLoggersLevel(c *rest.Context) {
	level := c.Param("level")
	if level == "" {
		return
	}

	l := logx.ParseLevel(level)

	if l == 0 {
		c.SendErrorJson(errorx.NewBadRequest("Error setting new logger level to " + level))
		return
	}

	logx.SetLevel(l)
	logx.Info("Changing Logger level to: ", l.String())
	c.String(http.StatusOK, "Logger level set to %s", l)
}

// swagger:route POST /logger/:name/:level update logger level
//
// Updates the log level of a specific logger within the service
//
// Responses:
//   200: Confirmation message
func updateLoggerLevel(c *rest.Context) {
	level, name := c.Param("level"), c.Param("name")
	if level == "" || name == "" {
		return
	}

	l := logx.ParseLevel(level)

	if l == 0 {
		c.SendErrorJson(errorx.NewBadRequest("Error setting new logger level to " + level))
		return
	}

	if logger := logx.Get(name); logger != nil {
		logger.SetLevel(l)
		c.String(http.StatusOK, "Logger level set to %s", l)
	} else {
		c.String(http.StatusNotFound, "Logger by the given name was not found")
	}

	logx.Infof("Changing level for logger=%s to=%s", name, l.String())
}

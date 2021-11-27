package endpoints

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-titan/info"
	"github.com/jucardi/go-titan/utils/paths"
)

var global = NewInfoHandler()

type InfoRouteHandler struct {
	routers map[string]*gin.Engine
}

func (h *InfoRouteHandler) RegisterEndpoint(router *gin.Engine) {
	router.GET("/info", h.getInfo)
}

func (h *InfoRouteHandler) AddRouter(address string, router *gin.Engine) {
	if h.routers == nil {
		h.routers = map[string]*gin.Engine{}
	}
	h.routers[address] = router
}

func (h *InfoRouteHandler) updateInfo() {
	var list []string
	for address, router := range h.routers {
		routes := router.Routes()
		l := streams.From(routes).
			Filter(func(i interface{}) bool {
				x := i.(gin.RouteInfo)
				return !strings.Contains(x.Path, "/__")
			}).
			Map(func(i interface{}) interface{} {
				x := i.(gin.RouteInfo)
				return fmt.Sprintf("%-7s %s", x.Method, paths.Combine(address, x.Path))
			}).
			ToArray().([]string)
		list = append(list, l...)
	}
	info.Set("routes", list)
}

// swagger:route GET /info health info
//
// Returns a information about this service
//
// Responses:
//   200: Info
func (h *InfoRouteHandler) getInfo(c *gin.Context) {
	h.updateInfo()
	c.IndentedJSON(200, info.Get())
}

func NewInfoHandler() *InfoRouteHandler {
	return &InfoRouteHandler{
		routers: map[string]*gin.Engine{},
	}
}

func InfoHandler() *InfoRouteHandler {
	return global
}

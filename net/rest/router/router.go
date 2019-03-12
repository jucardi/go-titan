package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type router struct {
	*routes
	base           *router
	group          *gin.RouterGroup
	groups         map[string]IRouter
	defaultSecured map[string]*routes
}

func (r *router) Base() IRouter {
	return r.base
}

func (r *router) Group(relativePath string, handlers ...HandlerFunc) IRouter {
	if router, ok := r.groups[relativePath]; ok {
		return router
	}

	group := wrap(r.group.Group(relativePath, convertHandlers(handlers)...))
	group.base = r.base
	r.groups[relativePath] = group
	return group
}

func (r *router) Version(apiVersion int) IRouter {
	version := fmt.Sprintf("v%d", apiVersion)
	return r.Group(version)
}

func wrap(r *gin.RouterGroup) *router {
	return &router{
		routes:         wrapRoutes(r),
		group:          r,
		groups:         map[string]IRouter{},
		defaultSecured: map[string]*routes{},
	}
}

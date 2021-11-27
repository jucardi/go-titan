package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-titan/net/rest"
)

type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes

	StaticFile(string, string) IRoutes
	Static(string, string) IRoutes
	StaticFS(string, http.FileSystem) IRoutes

	Before(...HandlerFunc) IRoutes
	After(...HandlerFunc) IRoutes
}

// IRouter defines the engine contracts
type IRouter interface {
	IRoutes

	// Base returns the router without a context path. You should only use this to add paths that are common in your services, such as /info or /metrics
	Base() IRouter

	// Version is a useful shortcut to create and/or retrieve the group containing routes defined for an API version.
	Version(apiVersion int) IRouter

	// Group creates a new engine group. You should add all the routes that have common middlwares or the same path prefix.
	// For example, all the routes that use a common middlware for authorization could be grouped.
	Group(relativePath string, handlers ...HandlerFunc) IRouter
}

// HandlerFunc defines the handler used by any middleware as return value.
type HandlerFunc func(*rest.Context)

func (h HandlerFunc) toGinHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		h(wrapContext(context))
	}
}

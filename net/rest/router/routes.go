package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type routes struct {
	before []HandlerFunc
	after  []HandlerFunc
	r      gin.IRoutes
}

func (r *routes) Before(handler ...HandlerFunc) IRoutes {
	for _, h := range handler {
		r.before = append(r.before, h)
	}
	return r
}

func (r *routes) After(handler ...HandlerFunc) IRoutes {
	for _, h := range handler {
		r.after = append(r.after, h)
	}
	return r
}

func (r *routes) Use(handlers ...HandlerFunc) IRoutes {
	r.r.Use(convertHandlers(handlers)...)
	return r
}

func (r *routes) Handle(method string, relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.Handle(method, relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.Any(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.GET(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.POST(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.DELETE(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.PATCH(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.PUT(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.OPTIONS(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	r.r.HEAD(relativePath, mergeHandlerGroups(r.before, handlers, r.after)...)
	return r
}

func (r *routes) StaticFile(relativePath, filepath string) IRoutes {
	r.r.StaticFile(relativePath, filepath)
	return r
}

func (r *routes) Static(relativePath, root string) IRoutes {
	r.r.Static(relativePath, root)
	return r
}

func (r *routes) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	r.r.StaticFS(relativePath, fs)
	return r
}

func wrapRoutes(r gin.IRoutes) *routes {
	return &routes{r: r}
}

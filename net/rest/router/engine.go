package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/rest/config"
	"github.com/jucardi/go-titan/net/rest/endpoints"
	"github.com/jucardi/go-titan/utils/paths"
)

type IEngine interface {
	IRouter

	// ServeHTTP conforms to the http.Handler interface.
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	// ContextPath returns the default context path for a service to listen to requests.
	// For example, for a user microservice, you me use "user" as the context path, and as a result
	// all routes added to the engine will be added using the "http://[url]:[port]/user/" prefix path.
	ContextPath() string

	// SetAdminAddress sets the address for the admin routes to be registered
	SetAdminAddress(addr ...string) IEngine

	// SetAddress sets the address where the router will listen for connections
	SetAddress(addr ...string) IEngine

	// Run attaches the engine to a http.Server and starts listening and serving HTTP requests using
	// the defined address with `SetAddress`. It is a shortcut for http.ListenAndServe(addr, engine)
	// Note: this method will block the calling goroutine indefinitely unless an error happens.
	Run() (err error)

	// RunAdmin asynchronously runs the admin engine. Returns an error if the admin port is undefined.
	RunAdmin() error

	// NoRoute adds handlers to use when a route is not found. It returns Not Found (404) by default.
	NoRoute(handlers ...HandlerFunc)
}

// Logger wrapper for Gin
type ginLogger struct {
}

func (l *ginLogger) Write(p []byte) (n int, err error) {
	logx.Trace(strings.TrimSpace(string(p)))
	return len(p), nil
}

func newLog() *ginLogger {
	return &ginLogger{}
}

func init() {
	onReloadConfig(config.Rest())
	config.AddReloadCallback(onReloadConfig)
}

func onReloadConfig(cfg *config.RestConfig) {
	l := newLog()
	gin.DefaultWriter = l
	if cfg.Verbose {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// Bare creates a new engine with no middlewares
func Bare(contextPath ...string) IEngine {
	return bare(contextPath...)
}

// BareFromConfig creates a new engine with the provided configuration and no middlewares
func BareFromConfig(cfg ...*config.RestConfig) IEngine {
	var c *config.RestConfig
	if len(cfg) > 0 && cfg[0] != nil {
		c = cfg[0]
	} else {
		c = config.Rest()
	}
	ret := bare(c.ContextPath)
	ret.config = *c
	return ret
}

// New creates a new Gin engine and applies the common middleware: Recovery, Logging, Metrics, Headers, Cid
func New(contextPath ...string) IEngine {
	router := Bare(contextPath...)
	UseCommonMiddleware(router)
	return router
}

// FromConfig creates a new Gin engine using the provided configuration and applies the common middleware:
// Recovery, Logging, Metrics, Headers, Cid
func FromConfig(cfg ...*config.RestConfig) IEngine {
	router := BareFromConfig(cfg...)
	UseCommonMiddleware(router)
	return router
}

func bare(contextPath ...string) *engine {
	r := &engine{
		engine: gin.New(),
		groups: map[string]gin.IRoutes{},
	}
	if len(contextPath) > 0 {
		r.contextPath = contextPath[0]
		r.router = wrap(r.engine.Group(r.ContextPath()))
	} else {
		r.router = wrap(&r.engine.RouterGroup)
	}

	r.router.base = wrap(&r.engine.RouterGroup)

	return r
}

type engine struct {
	*router
	engine      *gin.Engine
	adminEngine *gin.Engine
	contextPath string
	groups      map[string]gin.IRoutes
	config      config.RestConfig
	addr        []string
	adminAddr   []string
}

func (r *engine) SetAdminAddress(addr ...string) IEngine {
	r.adminAddr = addr
	return r
}

func (r *engine) SetAddress(addr ...string) IEngine {
	r.addr = addr
	return r
}

func (r *engine) ContextPath() string {
	return r.contextPath
}

func (r *engine) Run() error {
	r.runAdmin()
	routerAddr, _ := r.getAddresses()
	logx.Info("Listening on: ", routerAddr)
	endpoints.InfoHandler().AddRouter(paths.Combine(routerAddr...), r.engine)
	return r.engine.Run(routerAddr...)
}

func (r *engine) RunAdmin() error {
	if r.config.AdminPort == -1 {
		return errors.New("admin port is undefined")
	}
	r.runAdmin()
	return nil
}

func (r *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

func (r *engine) NoRoute(handlers ...HandlerFunc) {
	r.engine.NoRoute(mergeHandlerGroups(r.before, handlers, r.after)...)
}

func (r *engine) runAdmin() {
	if r.config.AdminPort == -1 {
		return
	}

	routerAddr, adminAddr := r.getAddresses()

	if len(adminAddr) == 0 || areAddressesEqual(routerAddr, adminAddr) {
		logx.Info("Admin endpoints registered on: ", routerAddr)
		registerAdminRoutes(r.engine)
		return
	}

	r.adminEngine = gin.New()
	registerAdminRoutes(r.adminEngine)
	endpoints.InfoHandler().AddRouter(paths.Combine(adminAddr...), r.adminEngine)

	go func() {
		logx.Info("Admin endpoints registered at: ", adminAddr)
		if err := r.adminEngine.Run(adminAddr...); err != nil {
			panic(fmt.Errorf("failed to start admin listener engine, %v", err))
		}
	}()
}

func (r *engine) getAddresses() (routerAddr []string, adminAddr []string) {
	routerAddr = r.addr
	adminAddr = r.adminAddr

	if len(routerAddr) == 0 {
		routerAddr = []string{fmt.Sprintf(":%d", r.config.HttpPort)}
	}

	if len(adminAddr) == 0 && r.config.AdminPort > 0 {
		adminAddr = []string{fmt.Sprintf(":%d", r.config.AdminPort)}
	}
	return
}

func registerAdminRoutes(r *gin.Engine) {
	endpoints.InfoHandler().RegisterEndpoint(r)
	endpoints.AddMetrics(r)
	endpoints.AddLogLevel(r)
}

func areAddressesEqual(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i, a1val := range a1 {
		if a1val != a2[i] {
			return false
		}
	}
	return true
}

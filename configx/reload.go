package configx

var (
	reloadHandlers []*callbackInfo
)

type callbackInfo struct {
	name string
	f    func(IConfig)
}

// Reload will trigger all registered OnReloadCallbacks. Useful to propagate a config change.
// Has no effect if a configuration has not been previously loaded.
func Reload() {
	if instance == nil {
		return
	}

	for _, h := range reloadHandlers {
		if h.name != "" {
			log().Debug("Triggering reload callback: ", h.name)
		}
		h.f(instance)
	}
}

// AddOnReloadCallback allows for a callback function to be registered that will be triggered on a config
// load or change. If registering a handler after the configuration has been loaded, the handler will be
// automatically executed once.
func AddOnReloadCallback(handler func(IConfig), name ...string) {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}

	reloadHandlers = append(reloadHandlers, &callbackInfo{name: n, f: handler})

	if instance != nil {
		handler(instance)
	}
}

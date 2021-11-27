package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Hook func() error

type hookInfo struct {
	name string
	fn   Hook
}

var (
	hooks []*hookInfo
)

// ListenForSignals for a TERM or INT signal.  Once the signal is caught all shutdown hooks will be
// executed allowing a graceful shutdown
func ListenForSignals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit

		log().Info("signal captured: ", sig.String())
		InvokeHooks()
		os.Exit(0)
	}()
}

// AddHook associates a no-arg func to be called when a signal is caught allowing for
// cleanup.
//
// Note: the function should not block for a period of time and hold up the shutdown
// of this app/service. All shutdown hooks must be independent of each other
// since they are executed concurrently for faster shutdown.
func AddHook(f Hook, name ...string) {
	hook := &hookInfo{
		fn: f,
	}
	if len(name) > 0 {
		hook.name = name[0]
	}
	hooks = append(hooks, hook)
}

func InvokeHooks() {
	var wg sync.WaitGroup

	wg.Add(len(hooks))

	for _, hook := range hooks {
		if hook.name != "" {
			log().Info("Executing hook: ", hook.name)
		}
		go func(hook *hookInfo) {
			defer wg.Done()
			if err := hook.fn(); err != nil {
				log().Warn("Shutdown hook exited with errors, ", err.Error())
			}
		}(hook)
	}

	wg.Wait()
}

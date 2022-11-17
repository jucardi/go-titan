package monitor

// IMonitor defines a monitor service which can run in the background of the main process and can have
// watchers attached to it
type IMonitor interface {
	// Start starts the monitor process, triggering `Run()` of all registered watchers on every tick
	Start()

	// StartAsync is similar to `Start` but the monitor process runs in a separate go routine so the
	// monitor process does not block the current thread.
	StartAsync()

	// Stop stops the monitor if it was previously started. Does nothing if the monitor is not currently
	// running
	Stop()

	// AddWatcher appends an implementation of `IWatcher` to the monitor so `Run` is triggered on every
	// tick of the monitor
	AddWatcher(watcher IWatcher)
}

// IWatcher defines a process that will be executed in the background to watch or track the state of a specific element.
// watchers can be registered in IMonitor instances
type IWatcher interface {
	// Run triggers the operation of a watcher
	Run()
}

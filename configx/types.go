package configx

const (
	FormatYaml ConfigFormat = "yaml"
	FormatJson ConfigFormat = "json"
)

type ConfigFormat string

type IConfig interface {
	// Value returns the value in the specified xPath, which can be a single key or a nested
	// path such as "rest.port". Returns nil if the value is not contained
	Value(xPath string, defaultVal ...interface{}) interface{}

	// String is the same as `Get` but automatically returns the value as `string`
	String(xPath string, defaultVal ...string) string

	// Boolean is the same as `Get` but automatically returns the value as `bool`
	Boolean(xPath string, defaultVal ...bool) bool

	// Int is the same as `Get` but automatically returns the value as `int`
	Int(xPath string, defaultVal ...int) int

	// Set sets the provided value to the provided path. If `makeParents` is provided and `true`, it will
	// automatically create parent nodes if they do not exist and initialize them as map[string]interface{}
	Set(xPath string, value interface{}, makeParents ...bool) error

	// AppName indicates the name of the application
	AppName() string

	// Load loads the config to the provided structure. Requires using the `cfg` tags which point to XPATH
	// in the loaded configuration so the values can be appended
	Load(obj interface{})

	// Hash returns the hash from raw data that was used to load the configuration contained by this instance
	Hash() string

	// MapToObj attempts to map the values present in the provider xPath to the provided struct
	MapToObj(xPath string, obj interface{}) error

	// GetMapped attempts to obtain a mapped object for the provided xPath from the internal cache.
	// Returns nil if it has not been previously mapped
	GetMapped(xPath string) interface{}
}

type RemotePullHandler func(url string) (data []byte, err error)

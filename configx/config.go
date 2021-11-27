package configx

import (
	"fmt"
	"sync"

	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/go-titan/utils/mapper"
	"github.com/jucardi/go-titan/utils/maps"
	"github.com/jucardi/go-titan/utils/reflectx"
)

const (
	KeyAppName = "app_name"
)

type Configuration struct {
	cfg           map[string]interface{}
	cache         map[string]interface{}
	onErrHandlers []func(error)
	mux           sync.Mutex
	source        string
	hash          string
	sourceFormat  ConfigFormat
}

// Value returns the value in the specified xPath, which can be a single key or a nested
// path such as "rest.port". Returns the provided defaultValue (or nil if none) if the value is not contained
func (b *Configuration) Value(xPath string, defaultVal ...interface{}) interface{} {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.get(xPath, defaultVal...)
}

// String is the same as `Value` but automatically returns the value as `string`
func (b *Configuration) String(xPath string, defaultVal ...string) string {
	return b.Value(xPath, stringx.GetOrDefault("", defaultVal...)).(string)
}

// Boolean is the same as `Value` but automatically returns the value as `bool`
func (b *Configuration) Boolean(xPath string, defaultVal ...bool) bool {
	def := false
	if len(defaultVal) > 0 {
		def = defaultVal[0]
	}
	return b.Value(xPath, def).(bool)
}

// Int is the same as `Value` but automatically returns the value as `int`
func (b *Configuration) Int(xPath string, defaultVal ...int) int {
	def := 0
	if len(defaultVal) > 0 {
		def = defaultVal[0]
	}
	return b.Value(xPath, def).(int)
}

func (b *Configuration) Set(xPath string, value interface{}, makeParents ...bool) error {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.set(xPath, value, makeParents...)
}

// AppName indicates the name of the application
func (b *Configuration) AppName() string {
	return b.Value(KeyAppName, "NO_NAME").(string)
}

func (b *Configuration) Load(obj interface{}) {

}

func (b *Configuration) Hash() string {
	return b.hash
}

// MapToObj converts a `map[string]interface{}` located in xPath to the provided structure
func (b *Configuration) MapToObj(xPath string, obj interface{}) error {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.addToCache(xPath, obj)
	val := b.get(xPath)
	if val == nil {
		log().Warn(fmt.Sprintf("value in path '%s' not found", xPath))
		return nil
	}

	mappingMode := mapper.MappingMode(b.sourceFormat)
	if err := mapper.Convert(obj, val, mappingMode); err != nil {
		return err
	}

	reflectx.Loader().Load(obj)

	newUpdate := map[string]interface{}{}

	if err := mapper.Convert(newUpdate, obj, mappingMode); err != nil {
		log().Warn("failed to map loaded values back to global config - ", err.Error())
	}

	if len(newUpdate) > 0 {
		if err := b.set(xPath, newUpdate); err != nil {
			log().Warn(fmt.Sprintf("failed to update global config with new mapped values - %v", err))
		}
	}

	return nil
}

// GetMapped returns the cached mapped structure for the provided xPath. Returns nil if none.
// Will return nil if `Convert` has not been invoked first
func (b *Configuration) GetMapped(xPath string) interface{} {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.cache == nil {
		b.cache = map[string]interface{}{}
	}
	if v, ok := b.cache[xPath]; ok {
		return v
	}
	return nil
}

func (b *Configuration) set(xPath string, value interface{}, makeParents ...bool) error {
	initParent := false
	if len(makeParents) > 0 {
		initParent = makeParents[0]
	}
	return maps.SetValue(b.cfg, xPath, value, initParent)
}

func (b *Configuration) get(xPath string, defaultVal ...interface{}) interface{} {
	var v interface{}
	if len(defaultVal) > 0 {
		v = defaultVal[0]
	}
	return maps.GetOrDefault(b.cfg, xPath, v)
}

func (b *Configuration) addToCache(key string, obj interface{}) {
	if b.cache == nil {
		b.cache = map[string]interface{}{}
	}
	b.cache[key] = obj
}

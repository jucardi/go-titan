package configx

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/go-titan/errors"
	"github.com/jucardi/go-titan/utils/maps"
	"gopkg.in/yaml.v3"
)

var (
	DefaultRemoteFreq = 1 * time.Minute

	instance IConfig

	hasher         = sha1.New()
	fieldCopyRegex = regexp.MustCompile(`^\${(.)*}$`)
	validFormats   = []ConfigFormat{FormatYaml, FormatJson}

	// Additional string aliases of the allowed formats
	formatAliases = map[string]ConfigFormat{
		"yml": FormatYaml,
	}

	loader *remoteLoader
)

// Get returns the loaded `IConfig` instance.
func Get() IConfig {
	if instance == nil {
		instance = &Configuration{
			cfg:   map[string]interface{}{},
			cache: map[string]interface{}{},
		}
	}
	return instance
}

// FromFile loads the configuration from a file into the global config instance.
func FromFile(path string) error {
	_, ext, err := validatePath(path)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("error reading file ", err.Error())
	}

	return load(ext, data)
}

// FromRemote loads the global configuration from the provided remote URL.
//
// On configuration loaded successfully, creates a watcher what will periodically retrieve the config
// from the same remote URL, and if changes ar found, triggers all reload callbacks registered to the
// manager. The pull frequency is defined by DefaultRemoteFreq. To disable recurrent pulls, set the
// value of DefaultRemoteFreq to 0 before calling this function.
//
// The expected format is determined by the extension ending of the URL (E.g. http://some/config.yml)
// to explicitly specify a expected config format, use FromRemoteWithFormat
//
//   {url}      -  The URL to be used to fetch the remote configuration
//   {handler}  -  If provided, it will override the default remote puller with the provided
//                 RemotePullHandler implementation.
//
func FromRemote(url string, handler ...RemotePullHandler) error {
	_, format, err := validatePath(url)
	if err != nil {
		return err
	}

	return FromRemoteWithFormat(url, format, handler...)
}

// FromRemoteWithFormat loads the global configuration from the provided remote URL.
//
// On configuration loaded successfully, creates a watcher what will periodically retrieve the config
// from the same remote URL, and if changes ar found, triggers all reload callbacks registered to the
// manager. The pull frequency is defined by DefaultRemoteFreq. To disable recurrent pulls, set the
// value of DefaultRemoteFreq to 0 before calling this function.
//
// The configuration will be loaded by attempting to deserialize the response body as the provided
// format (E.g yaml or json)
//
//   {url}      -  The URL to be used to fetch the remote configuration
//   {format}   -  Indicates the format the configuration is expected to be to properly deserialize it
//   {handler}  -  If provided, it will override the default remote puller with the provided
//                 RemotePullHandler implementation.
//
func FromRemoteWithFormat(url string, format ConfigFormat, handler ...RemotePullHandler) error {
	l, err := newRemoteLoader(url, format, handler...)
	if err != nil {
		return err
	}
	loader = l
	return loader.start()
}

func load(format ConfigFormat, data []byte) (err error) {
	hasher.Reset()
	hash := fmt.Sprintf("%s", hasher.Sum(data))
	if hash == Get().Hash() {
		log().Debug("no config changes detected")
		return
	}
	if Get().Hash() == "" {
		log().Info("loading configuration")
	} else {
		log().Info("detected config changes, reloading configuration")
	}
	cfg := &Configuration{
		cfg:    map[string]interface{}{},
		cache:  map[string]interface{}{},
		source: string(data),
		hash:   hash,
	}

	switch format {
	case FormatJson:
		cfg.sourceFormat = FormatJson
		err = json.Unmarshal(data, &cfg.cfg)
	case FormatYaml:
		cfg.sourceFormat = FormatYaml
		err = yaml.Unmarshal(data, cfg.cfg)
	default:
		err = errors.Format("no valid decoder found for '%s'", format)
	}
	if err == nil {
		processValuesCopy(cfg, cfg.cfg)
		instance = cfg
		Reload()
	}
	log().Debug("loaded configuration:\n" + string(data))
	return
}

func validatePath(path string) (filename string, format ConfigFormat, err error) {
	filename = filepath.Base(path)
	extStr := stringx.New(filepath.Ext(path)).ToLower().TrimLeft(".").S()

	if f, ok := formatAliases[extStr]; ok {
		format = f
	} else {
		format = ConfigFormat(extStr)
	}

	if filename == "" {
		return "", "", errors.New("filename cannot be empty")
	}

	err = validateFormat(format)
	return
}

func validateFormat(format ConfigFormat) error {
	if !streams.From(validFormats).Contains(format) {
		return errors.Format("unable to load configuration by the provided format '%s', supported formats are: %v", format, validFormats)
	}

	return nil
}

// Iterates over fields to determine if there are values that should be copied from a different
// XPATH location by using the ${copy_from_xpath} notation.
//
// E.g. given a configuration:
//
//   some_key: 'some value which can be a base type or a nested object'
//   some_other_key: '${some_key}'
//
// In that example 'some_other_key' is attempting to duplicate whatever value 'some_key'hast.
func processValuesCopy(orig *Configuration, cfg map[string]interface{}) {
	for k, v := range cfg {
		switch t := v.(type) {
		case string:
			if fieldCopyRegex.MatchString(t) {
				val := orig.Value(t[2 : len(t)-1])
				if val != nil {
					cfg[k] = processValue(orig, val)
				}
			}
		default:
			cfg[k] = processValue(orig, v)
		}
	}
}

// Custom value type handling.
//
// - If the value is a map[string]interface{} it will call the processValueCopy function to
//   determine if there are values that need to be copied from a different location of the global
//   configuration.
// - If the value is a map[interface{}]interface{} (common when deserializing from a YAML file)
//   it will attempt to convert it to a map[string]interface{} so then it can be handled by
//   processValueCopy
//
// Otherwise returns the same value as received in v
func processValue(orig *Configuration, v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		processValuesCopy(orig, t)
	case map[interface{}]interface{}:
		if val, err := maps.ConvertMap(t); err == nil {
			processValuesCopy(orig, val)
			return val
		}
	}
	return v
}

type remoteLoader struct {
	handler RemotePullHandler
	format  ConfigFormat
	url     string
	freq    time.Duration
}

func (r *remoteLoader) start() error {
	log().Debug("pulling remote config for the first time")
	err := r.triggerRemotePull()
	if err != nil {
		return err
	}

	if r.freq == 0 {
		log().Info("reload frequency is set to zero, will not start remote loader")
	} else {
		go r.ticker()
	}
	return nil
}

func (r *remoteLoader) ticker() {
	log().Debug("starting remote config watcher")
	t := time.NewTicker(r.freq)
	for {
		select {
		case <-t.C:
			if err := r.triggerRemotePull(); err != nil {
				log().Warn("failed tp fetch configuration - ", err.Error())
			}
		}
	}
}

func (r *remoteLoader) triggerRemotePull() error {
	if data, err := r.handler(r.url); err != nil {
		return err
	} else {
		return load(r.format, data)
	}
}

func newRemoteLoader(url string, format ConfigFormat, handler ...RemotePullHandler) (*remoteLoader, error) {
	if err := validateFormat(format); err != nil {
		return nil, err
	}

	ret := &remoteLoader{
		handler: defaultRemoteHandler,
		format:  format,
		url:     url,
		freq:    DefaultRemoteFreq,
	}

	if len(handler) > 0 && handler[0] != nil {
		ret.handler = handler[0]
	}

	return ret, nil
}

func defaultRemoteHandler(url string) (data []byte, err error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		err = errors.Format("failed to retrieve remote configuration - %v", err)
		return
	}
	if resp.StatusCode >= 400 {
		err = errors.Format("error status code (%d) received while retrieving remote configuration", resp.StatusCode)
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	return
}

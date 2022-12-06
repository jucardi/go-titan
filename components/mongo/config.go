package mongo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jucardi/go-titan/configx"
	"github.com/jucardi/go-titan/logx"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	// Name is a name for this connection to be used by the dependency manager.
	// This field is not use by mongo. Recommendation is to leave this blank if only connecting to a single mongo instance
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Host is the hostname where the database is located
	Host string `json:"host" yaml:"host" env:"MONGO_HOST"`

	// Port is the port where the database is listening to connections
	Port int `json:"port" yaml:"port" env:"MONGO_PORT"`

	// Database indicates the database name to connect to, where operations will be executed on
	Database string `json:"database" yaml:"database" env:"MONGO_DBNAME"`

	// AuthDB indicates the authentication database to use. If empty, the client will attempt to authenticate using the
	// 'Database' value instead
	AuthDB string `json:"auth_db" yaml:"auth_db" env:"MONGO_DBAUTH"`

	// Options is any additional options to be added to the connection string
	Options []string `json:"options,omitempty" yaml:"options,omitempty"`

	// Username is the username to authenticate to the database
	Username string `json:"username" yaml:"username" env:"MONGO_USERNAME"`

	// Password is the password to authenticate to the database
	Password string `json:"password" yaml:"password" env:"MONGO_PASSWORD"`

	// DialMaxRetries defines the maximum amount of retries to attempt when dialing to a db
	DialMaxRetries *int `json:"dial_max_retries" yaml:"dial_max_retries"`

	// DialRetryTimeout defines the timeout in milliseconds between retries when dialing to a db
	DialRetryTimeout *int64 `json:"dial_retry_timeout" yaml:"dial_retry_timeout"`

	// TlsCertLocation indicates the path for a PEM encoded cert for TLS. If supported, the client will attempt to
	// establish a TLS connection when this field is provided
	TlsCertLocation string `json:"tls_cert_path" yaml:"tls_cert_path" env:"DB_TLS_CERT_PATH"`

	// TlsSkipVerifyHost controls whether a client verifies the server's certificate chain and host name.
	// If TlsSkipVerifyHost is true, TLS accepts any certificate  presented by the server and any host name in that
	// certificate.  In this mode, TLS is susceptible to man-in-the-middle attacks. This should be used only for testing.
	TlsSkipVerifyHost bool `json:"tls_skip_verify_host" yaml:"tls_skip_verify_host"`

	// MigrationSource indicates the source where migration scripts are located
	MigrationSource string `json:"migration_source" yaml:"migration_source"`
}

func (c *Config) opts() *options.ClientOptions {
	creds := ""
	if c.Username != "" && c.Password != "" {
		creds = fmt.Sprintf("%s:%s@", url.QueryEscape(c.Username), url.QueryEscape(c.Password))
	}

	u := fmt.Sprintf("mongodb://%s%s:%d/%s%s", creds, c.Host, c.Port, c.Database, c.options())

	logx.Debug("mongodb connection string: ", u)
	ret := options.Client().ApplyURI(u)

	// TODO: Handler auth outside of the URL
	// TODO: Add TLS integration
	return ret
}

func (c *Config) url() string {
	auth := ""
	if c.Username != "" && c.Password != "" {
		auth = fmt.Sprintf("%s:%s@", c.Username, c.Password)
	}

	return fmt.Sprintf("mongodb://%s%s:%d/%s%s", auth, c.Host, c.Port, c.Database, c.options())
}

func (c *Config) dbName() string {
	if c.Database != "" {
		return c.Database
	}
	return "db"
}

func (c *Config) options() string {
	var opts []string
	if c.AuthDB != "" {
		opts = append(opts, "authSource="+c.AuthDB)
	}
	if c.TlsCertLocation == "" {
		opts = append(opts, "ssl=false")
	}
	opts = append(opts, c.Options...)
	if len(opts) == 0 {
		return ""
	}
	return "?" + strings.Join(opts, "&")
}

const (
	configKey  = "mongo"
	configName = "mongo-cfg"
)

var (
	singleton *Config
)

func init() {
	configx.AddOnReloadCallback(reloadCallback, configName)
}

// Rest returns rest configuration
func getConfig() *Config {
	if singleton == nil {
		reloadCallback(configx.Get())
	}
	return singleton
}

func reloadCallback(cfg configx.IConfig) {
	config := &Config{}

	logx.WithObj(
		cfg.MapToObj(configKey, config),
	).Fatal("unable to map service configuration")

	singleton = config
}

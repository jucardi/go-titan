package mongo

import (
	"fmt"

	"github.com/jucardi/go-titan/configx"
	"github.com/jucardi/go-titan/logx"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	// Name is a name for this connection to be used by the dependency manager.
	// This field is not use by mongo. Recommendation is to leave this blank if only connecting to a single mongo instance
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Host is the hostname where the database is located
	Host string `json:"host" yaml:"host"`

	// Port is the port where the database is listening to connections
	Port int `json:"port" yaml:"port"`

	// Database indicates the database name to connect to, where operations will be executed on
	Database string `json:"database" yaml:"database"`

	// AuthDB indicates the authentication database to use. If empty, the client will attempt to authenticate using the
	// 'Database' value instead
	AuthDB string `json:"auth_db" yaml:"auth_db"`

	// Options is any additional options to be added to the connection string
	Options string `json:"options,omitempty" yaml:"options,omitempty"`

	// Username is the username to authenticate to the database
	Username string `json:"username" yaml:"username"`

	// Password is the password to authenticate to the database
	Password string `json:"password" yaml:"password"`

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
}

func (c *Config) opts() *options.ClientOptions {
	creds := ""
	if c.Username != "" && c.Password != "" {
		creds = fmt.Sprintf("%s:%s@", c.Username, c.Password)
	}

	url := fmt.Sprintf("mongodb://%s%s:%d/%s%s", creds, c.Host, c.Port, c.Database, c.Options)
	logx.Debug("mongodb connection string: ", url)
	ret := options.Client().ApplyURI(url)

	// TODO: Handler auth outside of the URL
	// TODO: Add TLS integration
	return ret
}

func (c *Config) dbName() string {
	if c.Database != "" {
		return c.Database
	}
	return "db"
}

const (
	configKey  = "mongo"
	configName = "mongo-cfg"
)

var (
	singleton = &Config{}
)

func init() {
	configx.AddOnReloadCallback(func(cfg configx.IConfig) {
		config := &Config{}

		logx.WithObj(
			cfg.MapToObj(configKey, config),
		).Fatal("unable to map service configuration")

		singleton = config
	}, configName)
}

// Rest returns rest configuration
func getConfig() *Config {
	return singleton
}

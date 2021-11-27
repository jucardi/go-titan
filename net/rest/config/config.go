package config

const (
	EncodingAuto         Encoding = "auto"
	EncodingJson         Encoding = "json"
	EncodingIndentedJson Encoding = "indented-json"
	EncodingProto        Encoding = "proto"
)

type Encoding string

type RestConfig struct {
	// Port indicates the port where an API should be listening to
	HttpPort int `json:"http_port" yaml:"http_port" env:"TITAN_REST_HTTP_PORT"`

	// AdminPort is the port where the admin routes will be registered to.
	//  - If not set (0), the HttpPort will be used instead.
	//  - If set to -1, no admin endpoints will be registered
	//
	// NOTE: You should use a different port than HttpPort in production to avoid admin routes from
	// being exposed to the public
	AdminPort int `json:"admin_port" yaml:"admin_port" env:"TITAN_ADMIN_PORT"`

	// ContextPath indicates the context path to be used by the API
	ContextPath string `json:"context_path" yaml:"context_path" env:"TITAN_REST_CONTEXT_PATH"`

	// Response contains the configuration on how the response should be handled
	Response ResponseConfig `json:"response" yaml:"response"`

	// Reporting contains the configuration about how the middleware handles logging
	Reporting ReportingConfig `json:"reporting" yaml:"reporting"`

	// RequestLimitSize is the max byte size allowed in the request body. Zero means no limit. Default is 5Mib
	RequestLimitSize int64 `json:"request_limit_size" yaml:"request_limit_size" default:"5242880"`

	// Verbose enables verbose mode to the Gin router
	Verbose bool `json:"verbose" yaml:"verbose"`
}

type ResponseConfig struct {
	// Encoding indicates the responses will be encoded in controllers. If using the provided context functions
	// Send, SendOrErr, StatusOrErr, SendError, this will determine how the message will be encoded.
	// Does not apply for specific encoding functions such as Json, Protobuf, YAML, XML, etc
	Encoding Encoding `json:"mode" yaml:"mode" default:"auto"`

	// FallbackEncoding takes place if the default encoding fails to encode a message. Useful when using
	// `auto` which attempts to encode the response based on the incoming request Content-Type header, if
	// not provided the fallback mode will take place
	FallbackEncoding Encoding `json:"fallback_mode,omitempty" yaml:"fallback_mode,omitempty" default:"json"`

	// ErrorBodies indicates whether serializing errors and writing them to the response bodies should be enabled
	ErrorBodies bool `json:"error_bodies" yaml:"error_bodies" default:"true"`

	// ErrorStackTrace indicates whether stack traces should be sent to the client when an error occurs
	ErrorStackTrace bool `json:"error_stack_trace" yaml:"error_stack_trace"`
}

// ReportingConfig is the configuration for the middleware logging
type ReportingConfig struct {
	// MinStatus is the minimum HttpStatus code to do an `httputil.DumpRequest` to the logger by the
	// middleware. Any response which status code is equal or higher than the provided value will have
	// their http request dumped into the logger. The default is 500, to report all (5xx) server-side
	// errors and ignore all (4xx) client-side errors or (2xx) success statuses.
	MinStatus int `json:"min_status" yaml:"min_status" default:"500"`

	// Body indicates whether the request body should be dumped as well
	Body bool `json:"body" yaml:"body"`

	// StackTrace indicates whether the stack trace should be appended to the dump
	StackTrace bool `json:"stack_trace" yaml:"stack_trace"`
}

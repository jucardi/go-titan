package encoding

const (
	Json         = Encoding("json")
	IndentedJson = Encoding("indented-json")
	Gzip         = Encoding("gzip")
	Protobuf     = Encoding("proto")
	Auto         = Encoding("auto")
)

type Encoding string

package rest

// Header keys
const (
	// HeaderAuthorization is the standard authorization header used in rest
	HeaderAuthorization = "Authorization"

	// HeaderContentType is the standard Content-Type header used in rest
	HeaderContentType = "Content-Type"

	// HeaderResponseType is a custom header available in Jarvis to request responses in a specific encoding different
	// from the request content type.
	HeaderResponseType = "Response-Type"

	// HeaderContentEncoding is the standard Content-Encoder header user in rest
	HeaderContentEncoding = "Content-Encoder"
)

// Section for MIME Content Type values. Added a selection, could grow to add more or all.
//    ref: https://www.sitepoint.com/mime-types-complete-list/
const (
	// ContentTypeJson is the standard MIME type for json objects
	ContentTypeJson = "application/json"

	// ContentTypeProto is an non-standard yet commonly used MINE type for protobuf encoding
	ContentTypeProto = "application/x-protobuf"

	// ContentTypeYaml is an non-standard yet commonly used MINE type for YAML encoding
	ContentTypeYaml = "application/x-yaml"

	// ContentTypeText is the standard MIME type for javascript files
	ContentType = "application/javascript"

	// ContentTypeText is the standard MIME type for plain text
	ContentTypeText = "text/plain"

	// ContentTypeRichText is the standard MIME type for rich text
	ContentTypeRichText = "text/richtext"

	// ContentTypeCss is the standard MIME type for CSS files
	ContentTypeCss = "text/css"

	// ContentTypeHtml is the standard MIME type for HTML files
	ContentTypeHtml = "text/html"

	// ContentTypeBase64 is the standard MIME type for Base64 encoded bytes
	ContentTypeBase64 = "application/base64"

	// ContentTypeX509CaCert is the non-standard yet commonly used MIME type for X509 CA Certificates
	ContentTypeX509CaCert = "application/x-x509-ca-cert"

	// ContentTypeX509UserCert is the non-standard yet commonly used MIME type for X509 User Certificates
	ContentTypeX509UserCert = "application/x-x509-user-cert"

	// ContentTypePKIXCert is the standard MIME type for PKIX certificates
	ContentTypePKIXCert = "application/pkix-cert"

	// ContentTypePKCS7Mime is the content type for PKCS7 messages
	ContentTypePKCS7Mime = "application/pkcs7-mime"

	// ContentTypePKCS12 is the content type for PKCS12 cryptographic bundles
	ContentTypePKCS12 = "application/pkcs-12"

	// ContentTypeGzip is the content type for gzipped data
	ContentTypeGzip = "application/gzip"
)

package httpx

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

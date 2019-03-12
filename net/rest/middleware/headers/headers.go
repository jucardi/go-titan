package headers

import (
	"github.com/jucardi/go-titan/net/rest"
)

// Handler allows cross-origin requests to ease development
func Handler(c *rest.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, x-requested-by, *")
	c.Header("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
	c.Next()
}

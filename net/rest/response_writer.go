package rest

import (
	"github.com/gin-gonic/gin"
)

type respWriter struct {
	gin.ResponseWriter
}

func (r respWriter) WriteHeader(statusCode int) {
	if !r.ResponseWriter.Written() {
		r.ResponseWriter.WriteHeader(statusCode)
	}
}

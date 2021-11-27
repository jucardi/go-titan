package limits

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/net/errorx"
	"github.com/jucardi/go-titan/net/rest"
	"github.com/jucardi/go-titan/net/rest/config"
)

var maxSize int64 = 0

func init() {
	config.AddReloadCallback(func(config *config.RestConfig) {
		maxSize = config.RequestLimitSize
	})
}

type maxBytesReader struct {
	ctx        *rest.Context
	rdr        io.ReadCloser
	remaining  int64
	wasAborted bool
	sawEOF     bool
}

func (mbr *maxBytesReader) tooLarge() (n int, err error) {
	n, err = 0, fmt.Errorf("HTTP request too large")

	if !mbr.wasAborted {
		mbr.wasAborted = true
		ctx := mbr.ctx
		code := http.StatusRequestEntityTooLarge
		txt := http.StatusText(code)
		ctx.Header("connection", "close")
		ctx.AbortWithError(code, errorx.New(code, txt, txt))
		logx.Error(fmt.Sprintf("(%d) Request aborted, message too large", http.StatusRequestEntityTooLarge))
	}
	return
}

func (mbr *maxBytesReader) Read(p []byte) (n int, err error) {
	toRead := mbr.remaining
	if mbr.remaining == 0 {
		if mbr.sawEOF {
			return mbr.tooLarge()
		}
		// The underlying io.Reader may not return (0, io.EOF)
		// at EOF if the requested size is 0, so read 1 byte
		// instead. The io.Reader docs are a bit ambiguous
		// about the return value of Read when 0 bytes are
		// requested, and {bytes,strings}.Reader gets it wrong
		// too (it returns (0, nil) even at EOF).
		toRead = 1
	}
	if int64(len(p)) > toRead {
		p = p[:toRead]
	}
	n, err = mbr.rdr.Read(p)
	if err == io.EOF {
		mbr.sawEOF = true
	}
	if mbr.remaining == 0 {
		// If we had zero bytes to read remaining (but hadn't seen EOF)
		// and we get a byte here, that means we went over our limit.
		if n > 0 {
			return mbr.tooLarge()
		}
		return 0, err
	}
	mbr.remaining -= int64(n)
	if mbr.remaining < 0 {
		mbr.remaining = 0
	}
	return
}

func (mbr *maxBytesReader) Close() error {
	return mbr.rdr.Close()
}

// Handler is a middleware function that limits the size of request
// When a request is over the limit, the following will happen:
// * Error will be added to the context
// * Connection: close header will be set
// * Error 413 will be sent to the client (http.StatusRequestEntityTooLarge)
// * Current context will be aborted
func Handler(ctx *rest.Context) {
	if maxSize > 0 {
		ctx.Request.Body = &maxBytesReader{
			ctx:        ctx,
			rdr:        ctx.Request.Body,
			remaining:  maxSize,
			wasAborted: false,
			sawEOF:     false,
		}
	}
	ctx.Next()
}

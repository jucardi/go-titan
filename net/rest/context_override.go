package rest

func (c *Context) String(code int, format string, args ...interface{}) *Context {
	c.Context.String(code, format, args...)
	return c
}

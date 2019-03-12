package context

import "context"

type IContext interface {
	context.Context
	Get(key string) (val interface{}, exists bool)
	Set(key string, val interface{}) IContext
}

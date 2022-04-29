package ctx

import "context"

type ComponentBase struct {
	ctx context.Context
}

func (c *ComponentBase) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func New(ctx ...context.Context) *ComponentBase {
	if len(ctx) > 0 {
		return &ComponentBase{ctx: ctx[0]}
	}
	return &ComponentBase{}
}

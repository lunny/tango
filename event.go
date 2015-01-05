package tango

type Before interface{
	Before()
}

type After interface {
	After()
}

func Events() HandlerFunc {
	return func(ctx *Context) {
		action := ctx.Action()
		if action != nil {
			if b, ok := action.(Before); ok {
				b.Before()
			}
		}

		ctx.Next()

		if action != nil {
			if b, ok := action.(After); ok {
				b.After()
			}
		}
	}
}

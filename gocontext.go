package gocontext

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var contexts = new(sync.Map)

// GoWithContext start a new go routine with the given context
// and returns a context that is canceled with `fn` returns.
// It is also returns a cancel function that cancel the context
// provided within the Go routine. This cancel does not cancel
// the context you provide.
func GoWithContext(ctx context.Context, fn func()) (context.Context, func()) {

	childCtx, cancel := context.WithCancel(ctx)
	goCtx, cancelGo := context.WithCancel(ctx)

	go func() {
		id := goid()

		contexts.Store(id, childCtx)
		defer contexts.Delete(id)
		fn()
		cancelGo()
	}()

	return goCtx, cancel
}

// GoContext gives you the goroutine local context or context.Background
// if the goroutine has no context.
func GoContext() context.Context {
	if c, ok := contexts.Load(goid()); ok {
		return c.(context.Context)
	}

	return context.Background()
}

// based on https://github.com/golang/net/blob/master/http2/gotrack.go#L49-L68
var goroutineSpace = []byte("goroutine ")

func goid() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	n, err := strconv.ParseUint(string(b[:i]), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

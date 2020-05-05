package timelog

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
)

var WithOpenTracing = false

const tl_ctx_key = "tl_ctx_key"

// An information about an action and its internal actions.
// To create a context with it use the function Start().
type TlEntity struct {
	start    time.Time
	finish   time.Time
	message  interface{}
	otSpan   opentracing.Span
	parent   *TlEntity
	children []*TlEntity
}

func (e *TlEntity) finishAll(t time.Time) {
	if e.finish.IsZero() {
		if WithOpenTracing {
			e.otSpan.Finish()
		}
		e.finish = t
	}
	for _, child := range e.children {
		child.finishAll(t)
	}
}

// Starts measuring.
func Start(ctx context.Context, msg interface{}) context.Context {
	tl := &TlEntity{
		start:   time.Now(),
		message: msg,
	}

	if WithOpenTracing {
		tl.otSpan, ctx = opentracing.StartSpanFromContext(ctx, fmt.Sprintf("%s", msg))
	}

	if ctxTl := ctx.Value(tl_ctx_key); ctxTl != nil {
		tl.parent = ctxTl.(*TlEntity)
		ctxTl.(*TlEntity).children = append(ctxTl.(*TlEntity).children, tl)
	}

	return context.WithValue(ctx, tl_ctx_key, tl)
}

// Finishes measuring.
func Finish(ctx context.Context) context.Context {
	t := time.Now()
	var parent *TlEntity

	if ctxTl := ctx.Value(tl_ctx_key); ctxTl != nil {
		ctxTl.(*TlEntity).finishAll(t)
		parent = ctxTl.(*TlEntity).parent
	}

	if parent != nil {
		return context.WithValue(ctx, tl_ctx_key, parent)
	} else {
		return context.WithValue(ctx, tl_ctx_key, nil)
	}
}

// Returns a TimeLog entry from context.
// Returns nil if context does not containg *TlEntity.
func Get(ctx context.Context) *TlEntity {
	if ctxTl := ctx.Value(tl_ctx_key); ctxTl != nil {
		return ctxTl.(*TlEntity)
	}

	return nil
}

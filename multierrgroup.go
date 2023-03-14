// multierrgroup is a WaitGroup that works with functions that return errors.
package multierrgroup

import (
	"sync"

	"github.com/demosdemon/cpanic"
)

// Group is a collection of goroutines that can be waited on.
type Group struct {
	wg  sync.WaitGroup
	mu  sync.Mutex
	err manyErrors
}

// Wait blocks until all goroutines have finished and returns all errors, if any.
//
// For go1.20+, this will use the native error wrapping and return an error that
// implements `Unwrap() []error`. For go1.19 and below, this will use the Hashicorp
// multierror package.
func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err.done()
}

// AddError adds an error to the group. If the error is nil, it is ignored.
func (g *Group) AddError(err error) {
	if err == nil {
		return
	}

	g.mu.Lock()
	g.err.add(err)
	g.mu.Unlock()
}

// addPanic implements cpanic.Handler and adds the panic to the group.
func (g *Group) addPanic(p *cpanic.Panic) {
	g.AddError(p)
}

// Go runs the given function in a goroutine. If the function returns an error or
// panics, it is added to the group for later retrieval.
//
// It is not safe to call Go after Wait has been called.
func (g *Group) Go(fn func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		defer cpanic.Recover(g.addPanic)

		g.AddError(fn())
	}()
}

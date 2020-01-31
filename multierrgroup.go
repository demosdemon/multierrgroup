package multierrgroup

import (
	"sync"

	"github.com/demosdemon/cpanic"
	"github.com/hashicorp/go-multierror"
)

type Group struct {
	wg  sync.WaitGroup
	mu  sync.Mutex
	err *multierror.Error
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err.ErrorOrNil()
}

func (g *Group) AddError(err error) {
	if err == nil {
		return
	}

	g.mu.Lock()
	g.err = multierror.Append(g.err, err)
	g.mu.Unlock()
}

func (g *Group) addPanic(p *cpanic.Panic) {
	g.AddError(p)
}

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		defer cpanic.Recover(g.addPanic)

		g.AddError(fn())
	}()
}

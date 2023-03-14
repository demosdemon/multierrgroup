package multierrgroup

import (
	"sync"

	"github.com/demosdemon/cpanic"
)

type Group struct {
	wg  sync.WaitGroup
	mu  sync.Mutex
	err manyErrors
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err.done()
}

func (g *Group) AddError(err error) {
	if err == nil {
		return
	}

	g.mu.Lock()
	g.err.add(err)
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

package multierrgroup_test

import (
	"errors"
	"testing"

	"github.com/demosdemon/cpanic"
	"github.com/demosdemon/multierrgroup"
	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	var g multierrgroup.Group

	ch1 := make(chan func())
	ch2 := make(chan func())
	ch3 := make(chan func())
	ch4 := make(chan struct{})

	g.Go(func() error {
		f := <-ch1
		defer f()
		return nil
	})

	g.Go(func() error {
		f := <-ch2
		defer f()
		return errors.New("error")
	})

	g.Go(func() error {
		f := <-ch3
		defer f()
		panic("oh noes")
	})

	ch1 <- func() {
		t.Log("ch1")
		close(ch1)
		ch2 <- func() {
			t.Log("ch2")
			close(ch2)
			ch3 <- func() {
				t.Log("ch3")
				close(ch3)
				close(ch4)
			}
		}
	}

	<-ch4
	t.Log("ch4")
	err := g.Wait()
	assert.NotNil(t, err)

	merr := unwrap(err)
	assert.Len(t, merr, 2)

	stringErrors := make([]string, len(merr))
	for i, err := range merr {
		stringErrors[i] = err.Error()
	}

	assert.EqualValues(t, []string{"error", "panic: oh noes"}, stringErrors)

	_, ok := merr[1].(*cpanic.Panic)
	assert.True(t, ok)
}

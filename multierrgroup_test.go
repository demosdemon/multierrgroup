package multierrgroup_test

import (
	"errors"
	"testing"
	"time"

	"github.com/demosdemon/cpanic"
	"github.com/demosdemon/multierrgroup"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	var g multierrgroup.Group

	g.Go(func() error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	g.Go(func() error {
		time.Sleep(30 * time.Millisecond)
		return errors.New("error")
	})

	g.Go(func() error {
		time.Sleep(10 * time.Millisecond)
		panic("oh noes")
	})

	err := g.Wait()
	assert.NotNil(t, err)

	merr, ok := err.(*multierror.Error)
	assert.True(t, ok)
	assert.Len(t, merr.Errors, 2)

	stringErrors := make([]string, len(merr.Errors))
	for i, err := range merr.Errors {
		stringErrors[i] = err.Error()
	}

	assert.EqualValues(t, []string{"panic: oh noes", "error"}, stringErrors)

	_, ok = merr.Errors[0].(*cpanic.Panic)
	assert.True(t, ok)
}

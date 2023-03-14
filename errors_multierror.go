//go:build !go1.20

package multierrgroup

import "github.com/hashicorp/go-multierror"

type manyErrors struct{ *multierror.Error }

func (e *manyErrors) add(err error) {
	if err == nil {
		return
	}

	e.Error = multierror.Append(e.Error, err)
}

func (e *manyErrors) done() error {
	if e == nil {
		return nil
	}

	return e.ErrorOrNil()
}

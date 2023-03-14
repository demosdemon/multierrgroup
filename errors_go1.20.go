//go:build go1.20
// +build go1.20

package multierrgroup

import "errors"

type manyErrors []error

func (e *manyErrors) add(err error) {
	if err == nil {
		return
	}

	*e = append(*e, err)
}

func (e *manyErrors) done() error {
	if e == nil {
		return nil
	}

	return errors.Join((*e)...)
}

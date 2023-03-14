//go:build !go1.20

package multierrgroup_test

import "github.com/hashicorp/go-multierror"

func unwrap(err error) []error {
	return err.(*multierror.Error).Errors
}

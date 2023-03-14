//go:build go1.20

package multierrgroup_test

func unwrap(err error) []error {
	return err.(interface {
		Unwrap() []error
	}).Unwrap()
}

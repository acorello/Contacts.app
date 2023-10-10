package http

import "testing"

func TestNameRegex(t *testing.T) {
	val := "Joe!"
	err := nameRegEx.Validate(val)
	if err == nil {
		t.Errorf("%q should have failed", val)
	}
}

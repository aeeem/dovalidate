package validate_test

import (
	validate "dovalidate"
	"testing"
)

type TestStruct struct {
	A string `validate:"required,numeric" json:"otp"`
	B string `validate:"required"`
}

func TestValidating(t *testing.T) {
	validate.RegisterCustomTags()
	a := TestStruct{}
	if err := validate.Struct(a); err != nil {
		t.Log(err)
		return
	}
	t.Log(a)
}

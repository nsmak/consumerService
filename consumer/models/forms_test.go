package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegFromValidateFail(t *testing.T) {
	for _, tst := range []RegFrom{
		{
			Email: "",
			Pass1: "",
			Pass2: "",
		},
		{
			Email: "test",
			Pass1: "pass",
			Pass2: "pass",
		},
		{
			Email: "test@test.com",
			Pass1: "pass",
			Pass2: "",
		},
		{
			Email: "test@test.com",
			Pass1: "",
			Pass2: "pass",
		},
		{
			Email: "test@test.com",
			Pass1: "pass",
			Pass2: "pass1",
		},
		{
			Email: "test@test.com",
			Pass1: "pass1",
			Pass2: "pass",
		},
	} {
		err := tst.Validate()
		require.Error(t, err, "case: ", tst)
	}
}

func TestRegFromValidateSuccess(t *testing.T) {
	form := RegFrom{
		Email: "test@test.com",
		Pass1: "123456",
		Pass2: "123456",
	}
	err := form.Validate()
	require.NoError(t, err)
}

func TestAuthFormValidateFail(t *testing.T) {
	for _, tst := range []AuthForm{
		{
			Email: "",
			Pass:  "",
		},
		{
			Email: "",
			Pass:  "123456",
		},
		{
			Email: "test@mail.ru",
			Pass:  "",
		},
		{
			Email: "test",
			Pass:  "123456",
		},
	} {
		err := tst.Validate()
		require.Error(t, err, "case: ", tst)
	}
}

func TestAuthFormValidateSuccess(t *testing.T) {
	form := AuthForm{
		Email: "test@test.com",
		Pass:  "123456",
	}
	err := form.Validate()
	require.NoError(t, err)
}

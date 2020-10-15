package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainError_Error(t *testing.T) {
	t.Run("without dependent error", func(t *testing.T) {
		err := &domainError{IsUserError: true, Message: "err"}
		require.Equal(t, "err", err.Error())
		require.Equal(t, err.IsUserError, err.UserError())
	})

	t.Run("with dependent error", func(t *testing.T) {
		err1 := &domainError{IsUserError: true, Message: "err1"}
		err2 := domainError{Message: "err2", Err: err1}
		require.Equal(t, "err2 --> err1", err2.Error())
		require.NotEqual(t, err2.IsUserError, err2.UserError())
	})

	t.Run("when dependent error is not \"consumer.Error\"", func(t *testing.T) {
		err := &domainError{IsUserError: true, Message: "err", Err: errors.New("error")}
		require.Equal(t, "err --> error", err.Error())
		require.Equal(t, err.IsUserError, err.UserError())
	})
}

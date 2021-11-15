package jsonbind

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBindPathError(t *testing.T) {
	err := BindPathError("abc")
	var _e BindPathError
	require.True(t, errors.As(err, &_e))
}

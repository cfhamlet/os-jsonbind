package jsonbind

import "errors"

var (
	ErrCast         = errors.New("cast error")
	ErrPreallocate  = errors.New("preallocate error")
	ErrNotSupported = errors.New("not supported")
)

type BindPathError string

func (e BindPathError) Error() string {
	return string(e)
}

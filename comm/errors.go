package comm

import (
	"errors"
)

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrNotFound         = errors.New("not found")
	ErrAleadyExists     = errors.New("already exist")
	ErrPermissionDenied = errors.New("permission denied")
	ErrAborted          = errors.New("aborted")
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrInternal         = errors.New("internal")
)

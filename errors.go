package nbattle

import "errors"

var (
	ErrObjectNotFound       = errors.New("object not found")
	ErrUnexpectedObjectType = errors.New("unexpected object type")
)

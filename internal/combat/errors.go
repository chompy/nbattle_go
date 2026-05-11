package combat

import "errors"

var (
	ErrObjectNotFound          = errors.New("object not found")
	ErrUnexpectedObjectType    = errors.New("unexpected object type")
	ErrUnserializableValueType = errors.New("value cannot be serialized")
	ErrDeserializeWrongType    = errors.New("serialized data is not for this object type")
	ErrNilObject               = errors.New("object is nil")
)

package event

import "errors"

var (
	ErrUnserializableValueType = errors.New("value cannot be serialized")
	ErrDeserializeWrongType    = errors.New("serialized data is not for this object type")
)

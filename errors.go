package nbattle

import (
	"errors"
	"log"
	"strings"
)

var (
	ErrObjectNotFound          = errors.New("object not found")
	ErrUnexpectedObjectType    = errors.New("unexpected object type")
	ErrUnserializableValueType = errors.New("value cannot be serialized")
	ErrDeserializeWrongType    = errors.New("serialized data is not for this object type")
	ErrNilObject               = errors.New("object is nil")
	ErrUnexpectedLuaFuncReturn = errors.New("unexpected return value/type from lua function call")
)

func logLuaFuncCallError(err error, funcName string) {
	if !isLuaNotFoundErr(err) {
		log.Printf("WARNING: Error during Lua function call to %s: %s", funcName, err.Error())
	}
}

func isLuaNotFoundErr(err error) bool {
	errStr := err.Error()
	return strings.HasPrefix(errStr, "function") && strings.HasSuffix(errStr, "not found")
}

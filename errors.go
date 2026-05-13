package nbattle

import (
	"errors"
	"log/slog"
	"strings"
)

var (
	ErrObjectNotFound          = errors.New("object not found")
	ErrUnexpectedObjectType    = errors.New("unexpected object type")
	ErrNilObject               = errors.New("object is nil")
	ErrUnexpectedLuaFuncReturn = errors.New("unexpected return value/type from lua function call")
)

func logLuaFuncCallError(logger *slog.Logger, err error, funcName string) {
	if !isLuaNotFoundErr(err) {
		logger.Warn("Error during Lua function call.", "function", funcName, "error", err)
	}
}

func isLuaNotFoundErr(err error) bool {
	errStr := err.Error()
	return strings.HasPrefix(errStr, "function") && strings.HasSuffix(errStr, "not found")
}

package lua

import (
	"errors"
	"log"
)

var (
	ErrUnexpectedLuaFuncReturn = errors.New("unexpected return value/type from lua function call")
)

func logLuaFuncCallError(err error, funcName string) {
	log.Printf("WARNING: Error during Lua function call to %s: %s", funcName, err.Error())
}

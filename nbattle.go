package nbattle

import (
	"io"

	"github.com/chompy/nbattle_go/internal/combat"
	"github.com/chompy/nbattle_go/internal/lua"
)

// New creates a new NBattle context.
func New() *combat.Context {
	return combat.New()
}

// NewLuaEffect creates a new effect definition from a Lua script.
func NewLuaEffect(ctx *combat.Context, script io.Reader) (*combat.EffectDef, error) {
	return lua.NewEffect(ctx, script)
}

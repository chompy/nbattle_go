package nbattle

import (
	"io"

	"github.com/chompy/nbattle_go/internal/combat"
	"github.com/chompy/nbattle_go/internal/lua"
)

type Context *combat.Context
type StatDef *combat.StatDef
type EffectDef *combat.EffectDef
type TriggerDef *combat.TriggerDef
type Stat *combat.Stat
type Combatant *combat.Combatant

// New creates a new NBattle context.
func New() Context {
	return combat.New()
}

// NewLuaEffect creates a new effect definition from a Lua script.
func NewLuaEffect(ctx *combat.Context, script io.Reader) (EffectDef, error) {
	return lua.NewEffect(ctx, script)
}

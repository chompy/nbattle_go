package lua

import (
	"fmt"

	nbattle "github.com/chompy/nbattle_go"
)

func StatDefToLua(statDef *nbattle.StatDef) map[string]any {
	if statDef == nil {
		return ErrorToLua(nbattle.ErrNilObject)
	}
	return map[string]any{
		"type": statDef.GetType(),
		"name": statDef.GetName(),
		"min":  statDef.GetMin(),
		"max":  statDef.GetMax(),
	}
}

func StatToLua(ctx *nbattle.Context, stat *nbattle.Stat) map[string]any {
	return map[string]any{
		"type":    stat.GetType(),
		"def":     StatDefToLua(stat.GetDef()),
		"getBase": stat.GetBase,
		"setBase": func(value float64) {
			stat.SetBase(int(value))
		},
		"addBase": func(value float64) {
			stat.AddBase(int(value))
		},
		"getValue": stat.GetValue,
		"setMod": func(source any, value float64) {
			sourceObj, err := ObjectFromLua(ctx, source)
			if err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("stat.%d.setMod", stat.GetID()))
				return
			}
			stat.SetMod(sourceObj, int(value))
		},
		"get": stat.GetValue,
		"set": func(value float64) {
			stat.SetBase(int(value))
		},
		"add": func(value float64) {
			stat.AddBase(int(value))
		},
		"getCombatant": func() map[string]any {
			combatant, _ := ctx.GetCombatantWithStat(stat)
			return ObjectToLua(ctx, combatant)
		},
	}
}

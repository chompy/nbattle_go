package lua

import (
	nbattle "github.com/chompy/nbattle_go"
)

func statDefToLua(statDef *nbattle.StatDef) map[string]any {
	if statDef == nil {
		return errorToLua(nbattle.ErrNilObject)
	}
	return map[string]any{
		"type": statDef.GetType(),
		"name": statDef.GetName(),
		"min":  statDef.GetMin(),
		"max":  statDef.GetMax(),
	}
}

func statToLua(ctx *nbattle.Context, stat *nbattle.Stat) map[string]any {
	return map[string]any{
		"def":     statDefToLua(stat.GetDef()),
		"getBase": stat.GetBase,
		"setBase": func(value float64) {
			stat.SetBase(int(value))
		},
		"addBase": func(value float64) {
			stat.AddBase(int(value))
		},
		"getValue": stat.GetValue,
		"setMod": func(source any, value int) {
			sourceObj, err := objectFromLua(ctx, source)
			if err != nil {
				logLuaFuncCallError(err, "stat.setMod")
				return
			}
			stat.SetMod(sourceObj, value)
		},
		"get": stat.GetValue,
		"set": func(value float64) {
			stat.SetBase(int(value))
		},
		"add": func(value float64) {
			stat.AddBase(int(value))
		},
		"subtract": func(value float64) {
			stat.SubBase(int(value))
		},
		"sub": func(value float64) {
			stat.SubBase(int(value))
		},
		"getCombatant": func() map[string]any {
			combatant, _ := ctx.GetCombatantWithStat(stat)
			return combatantToLua(ctx, combatant)
		},
	}
}

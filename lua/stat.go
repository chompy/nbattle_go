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
		"subBase": func(value float64) {
			stat.SubBase(int(value))
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
		"getCombatant": func() map[string]any {
			combatant, err := ctx.GetCombatantWithStat(stat)
			if err != nil {
				return errorToLua(err)
			}
			return combatantToLua(ctx, combatant)
		},
	}
}

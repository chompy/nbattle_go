package lua

import (
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
		"def":     StatDefToLua(stat.GetDef()),
		"getBase": stat.GetBase,
		"setBase": func(value float64) {
			stat.SetBase(int(value))
		},
		"addBase": func(value float64) {
			stat.AddBase(int(value))
		},
		"getValue": stat.GetValue,
		"setMod": func(source any, value int) {
			sourceObj, err := ObjectFromLua(ctx, source)
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
		"sub": func(value float64) {
			stat.SubBase(int(value))
		},
		"getCombatant": func() map[string]any {
			combatant, _ := ctx.GetCombatantWithStat(stat)
			return CombatantToLua(ctx, combatant)
		},
	}
}

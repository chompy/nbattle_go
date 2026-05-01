package lua

import (
	"log"

	nbattle "github.com/chompy/nbattle_go"
)

func ContextToLua(ctx *nbattle.Context) map[string]any {
	return map[string]any{
		"getTick": ctx.GetTick,
		"getCombatants": func() []map[string]any {
			out := make([]map[string]any, 0)
			for _, combatant := range ctx.Combatants() {
				out = append(out, CombatantToLua(combatant))
			}
			return out
		},
	}
}

func StatDefToLua(stat *nbattle.StatDef) map[string]any {
	return map[string]any{
		"name": stat.GetName(),
		"min":  stat.GetMin(),
		"max":  stat.GetMax(),
	}
}

func StatToLua(stat *nbattle.Stat) map[string]any {
	return map[string]any{
		"def":     StatDefToLua(stat.GetDef()),
		"getBase": stat.GetBase,
		"setBase": stat.SetBase,
		"get":     stat.GetValue,
	}
}

func CombatantToLua(combatant *nbattle.Combatant) map[string]any {
	stats := make(map[string]any)
	for _, stat := range combatant.Stats() {
		stats[stat.GetDef().GetName()] = StatToLua(stat)
	}
	return map[string]any{
		"id":    combatant.GetID(),
		"stats": stats,
		"addEffect": func(name string, sourceObj any) {
			sourceID := 0
			source, ok := sourceObj.(map[string]any)
			if ok {
				sourceID = source["id"].(int)
			}
			if err := combatant.AddEffect(name, sourceID); err != nil {
				log.Println("WARNING: Error during Lua function call to combatant.addEffect:", err)
			}
		},
	}
}

func EffectContextToLua(ctx *nbattle.EffectCtx) map[string]any {
	return map[string]any{
		"source": CombatantToLua(ctx.Source),
		"target": CombatantToLua(ctx.Target),
	}
}

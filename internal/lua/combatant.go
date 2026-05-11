package lua

import (
	"fmt"

	nbattle "github.com/chompy/nbattle_go/internal/combat"
)

func combatantToLua(ctx *nbattle.Context, combatant *nbattle.Combatant) map[string]any {
	if combatant == nil {
		return map[string]any{
			"type":  nbattle.ObjectTypeUnknown,
			"error": nbattle.ErrNilObject,
		}
	}
	return map[string]any{
		"id":   combatant.GetID(),
		"type": combatant.GetType(),
		"getStat": func(statDef any) map[string]any {
			stat, err := combatant.GetStat(statDef)
			if err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("combatant.%d.getStat", combatant.GetID()))
				return errorToLua(err)
			}
			return statToLua(ctx, stat)
		},
		"setEffect": func(effectDef any, potency int, sourceObj any) {
			if err := combatant.SetEffect(effectDef, potency, sourceObj); err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("combatant.%d.setEffect", combatant.GetID()))
			}
		},
		"hasEffect": combatant.HasEffect,
		"removeEffect": func(effectDef any) {
			if err := combatant.SetEffect(effectDef, 0, nil); err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("combatant.%d.removeEffect", combatant.GetID()))
			}
		},
		"setFlag": combatant.SetFlag,
		"hasFlag": combatant.HasFlag,
	}
}

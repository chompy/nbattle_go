package lua

import (
	"fmt"

	nbattle "github.com/chompy/nbattle_go"
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
		"getStat": func(statDefName string) map[string]any {
			return statToLua(ctx, combatant.GetStat(statDefName))
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
	}
}

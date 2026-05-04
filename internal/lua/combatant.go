package lua

import (
	"fmt"

	nbattle "github.com/chompy/nbattle_go"
)

func CombatantToLua(ctx *nbattle.Context, combatant *nbattle.Combatant) map[string]any {
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
			return StatToLua(ctx, combatant.GetStat(statDefName))
		},
		"setEffect": func(effectDef any, potency int, sourceObj any) {
			if err := combatant.SetEffect(effectDef, potency, sourceObj); err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("combatant.%d.addEffect", combatant.GetID()))
			}
		},
		"hasEffect": combatant.HasEffect,
	}
}

func CombatantFromLua(combatant any, ctx *nbattle.Context) (*nbattle.Combatant, error) {
	switch combatant := combatant.(type) {
	case map[string]any:
		id, ok := combatant["id"].(int)
		if !ok {
			return nil, nbattle.ErrUnexpectedObjectType
		}
		return ctx.GetCombatantByID(id)
	case int:
		return ctx.GetCombatantByID(combatant)
	case float32:
		return ctx.GetCombatantByID(int(combatant))
	case float64:
		return ctx.GetCombatantByID(int(combatant))
	default:
		return nil, nbattle.ErrUnexpectedObjectType
	}
}

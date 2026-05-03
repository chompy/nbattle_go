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
	stats := make(map[string]any)
	for _, stat := range combatant.GetStats() {
		stats[stat.GetDef().GetName()] = StatToLua(ctx, stat)
	}
	return map[string]any{
		"id":   combatant.GetID(),
		"type": combatant.GetType(),
		"stat": stats,
		"addEffect": func(name string, sourceObj any) {
			sourceID := 0
			source, ok := sourceObj.(map[string]any)
			if ok {
				sourceID = source["id"].(int)
			}
			if err := combatant.AddEffect(name, sourceID); err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("combatant.%d.addEffect", combatant.GetID()))
			}
		},
		"removeEffect": func(name string) {
			if err := combatant.RemoveEffect(name); err != nil {
				logLuaFuncCallError(err, fmt.Sprintf("combatant.%d.removeEffect", combatant.GetID()))
			}
		},
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

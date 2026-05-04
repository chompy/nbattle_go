package lua

import (
	"log"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/event"
)

func luaGlobals(ctx *nbattle.Context) map[string]any {
	return map[string]any{
		"TICK":                event.TickEvent,
		"COMBATANT_UPDATE":    event.CombatantUpdateEvent,
		"COMBATANT_STAT_BASE": event.CombatantStatBaseEvent,
		"COMBATANT_STAT_MOD":  event.CombatantStatModEvent,
		"COMBATANT_EFFECT":    event.CombatantEffectEvent,
		"STAT_DEF":            nbattle.ObjectTypeStatDef,
		"EFFECT_DEF":          nbattle.ObjectTypeEffectDef,
		"COMBATANT":           nbattle.ObjectTypeCombatant,
		"ctx":                 contextToLua(ctx),
	}
}

func contextToLua(ctx *nbattle.Context) map[string]any {
	return map[string]any{
		"getTick": ctx.GetTick,
		"getCombatants": func() []map[string]any {
			out := make([]map[string]any, 0)
			for _, combatant := range ctx.GetCombatants() {
				out = append(out, combatantToLua(ctx, combatant))
			}
			return out
		},
		"getObject": func(obj any) {
			objectToLua(ctx, obj)
		},
	}
}

func errorToLua(err error) map[string]any {
	log.Println("WARNING: Error during Lua call:", err)
	return map[string]any{
		"type":  nbattle.ObjectTypeError,
		"error": err.Error(),
	}
}

func objectIDFromLua(object any) (int, error) {
	switch object := object.(type) {
	case map[string]any:
		id, ok := object["id"].(int)
		if !ok {
			return 0, nbattle.ErrUnexpectedObjectType
		}
		return id, nil

	case float64:
		return int(object), nil

	case float32:
		return int(object), nil

	case int:
		return object, nil

	}
	return 0, nbattle.ErrUnexpectedObjectType
}

func objectFromLua(ctx *nbattle.Context, object any) (nbattle.Object, error) {
	id, err := objectIDFromLua(object)
	if err != nil {
		return nil, err
	}
	obj := ctx.GetObjectByID(id)
	if obj == nil {
		return nil, nbattle.ErrObjectNotFound
	}
	return obj, nil
}

func objectToLua(ctx *nbattle.Context, object any) map[string]any {
	switch object := object.(type) {
	case *nbattle.StatDef:
		return statDefToLua(object)
	case *nbattle.Stat:
		return statToLua(ctx, object)
	case *nbattle.Combatant:
		return combatantToLua(ctx, object)
	}
	return map[string]any{
		"type":  "unknown",
		"error": nbattle.ErrUnexpectedObjectType,
	}
}

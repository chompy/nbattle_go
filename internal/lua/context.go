package lua

import (
	"log"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/internal/event"
)

func LuaGlobals(ctx *nbattle.Context) map[string]any {
	return map[string]any{
		"TICK":                    event.TickEvent,
		"STAT_BASE":               event.StatBaseEvent,
		"STAT_MOD":                event.StatModEvent,
		"NEW_COMBATANT":           event.NewCombatantEvent,
		"ADD_COMBATANT_EFFECT":    event.AddCombatantEffectEvent,
		"REMOVE_COMBATANT_EFFECT": event.RemoveCombatantEffectEvent,
		"STAT_DEF":                nbattle.ObjectTypeStatDef,
		"STAT":                    nbattle.ObjectTypeStat,
		"EFFECT_DEF":              nbattle.ObjectTypeEffectDef,
		"COMBATANT":               nbattle.ObjectTypeCombatant,
		"ctx":                     ContextToLua(ctx),
	}
}

func ContextToLua(ctx *nbattle.Context) map[string]any {
	return map[string]any{
		"getTick": ctx.GetTick,
		"getCombatants": func() []map[string]any {
			out := make([]map[string]any, 0)
			for _, combatant := range ctx.GetCombatants() {
				out = append(out, CombatantToLua(ctx, combatant))
			}
			return out
		},
		"getObject": func(obj any) {},
	}
}

func ErrorToLua(err error) map[string]any {
	log.Println("WARNING: Error during Lua call:", err)
	return map[string]any{
		"type":  nbattle.ObjectTypeError,
		"error": err.Error(),
	}
}

func ObjectIDFromLua(object any) (int, error) {
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

func ObjectFromLua(ctx *nbattle.Context, object any) (nbattle.Object, error) {
	id, err := ObjectIDFromLua(object)
	if err != nil {
		return nil, err
	}
	obj := ctx.GetObjectByID(id)
	if obj == nil {
		return nil, nbattle.ErrObjectNotFound
	}
	return obj, nil
}

func ObjectToLua(ctx *nbattle.Context, object any) map[string]any {
	switch object := object.(type) {
	case *nbattle.StatDef:
		return StatDefToLua(object)
	case *nbattle.Stat:
		return StatToLua(ctx, object)
	case *nbattle.Combatant:
		return CombatantToLua(ctx, object)
	}
	return map[string]any{
		"type":  "unknown",
		"error": nbattle.ErrUnexpectedObjectType,
	}
}

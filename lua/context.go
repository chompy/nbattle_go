package lua

import (
	"log"
	"strings"

	nbattle "github.com/chompy/nbattle_go"
)

func luaGlobals(ctx *nbattle.Context) map[string]any {

	globals := make(map[string]any)
	globals["ctx"] = contextToLua(ctx)
	for _, statDef := range ctx.GetStatDefs() {
		globals["STAT_"+strings.ToUpper(statDef.GetName())] = statDefToLua(statDef)
	}
	for name, value := range ctx.GetFlags() {
		globals["FLAG_"+strings.ToUpper(name)] = value
	}
	return globals
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
		"getObject": func(obj any) map[string]any {
			return objectToLua(ctx, obj)
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
	obj, err := ctx.GetObjectByID(id)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func objectToLua(ctx *nbattle.Context, object any) map[string]any {
	switch object := object.(type) {
	case *nbattle.StatDef:
		return statDefToLua(object)

	case *nbattle.EffectDef:
		return map[string]any{
			"id":   object.GetID(),
			"type": object.GetType(),
			"name": object.GetName(),
		}

	case *nbattle.TriggerDef:
		return map[string]any{
			"id":   object.GetID(),
			"type": object.GetType(),
			"name": object.GetName(),
		}

	case *nbattle.Stat:
		return statToLua(ctx, object)

	case *nbattle.Combatant:
		return combatantToLua(ctx, object)

	default:
		return errorToLua(nbattle.ErrUnexpectedObjectType)
	}
}

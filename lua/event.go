package lua

import (
	"log"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/event"
)

func eventToLua(ctx *nbattle.Context, evt event.Event) map[string]any {
	out := make(map[string]any)
	out["type"] = evt.Type()
	switch evt := evt.(type) {
	case *event.Tick:
		out["tick"] = evt.Tick

	case *event.CombatantUpdate:
		out["combatant"] = nil
		out["active"] = evt.Active
		combatant, err := ctx.GetCombatantByID(evt.CombatantID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		out["combatant"] = combatantToLua(ctx, combatant)

	case *event.CombatantStatBase:
		out["value"] = evt.Value
		out["setValue"] = func(value float64) {
			evt.Value = int(value)
		}
		combatant, err := ctx.GetCombatantByID(evt.CombatantID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		out["combatant"] = combatantToLua(ctx, combatant)
		out["stat"] = statToLua(ctx, combatant.GetStat(evt.StatDefID))

	case *event.CombatantStatMod:
		out["modValue"] = evt.ModValue
		out["setModValue"] = func(value float64) {
			evt.ModValue = int(value)
		}
		out["source"] = objectToLua(ctx, ctx.GetObjectByID(evt.SourceID))

		combatant, err := ctx.GetCombatantByID(evt.CombatantID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		out["combatant"] = combatant
		out["stat"] = statToLua(ctx, combatant.GetStat(evt.StatDefID))

	case *event.CombatantEffect:
		effectDef, err := ctx.GetEffectDefByID(evt.EffectDefID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		target, err := ctx.GetCombatantByID(evt.TargetID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		out["target"] = combatantToLua(ctx, target)
		out["effect"] = effectDef.GetName()
		out["potency"] = evt.Potency
		out["source"] = objectToLua(ctx, ctx.GetObjectByID(evt.SourceID))
	}

	return out
}

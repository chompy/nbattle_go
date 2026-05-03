package lua

import (
	"log"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/internal/event"
)

func EventToLua(ctx *nbattle.Context, evt event.Event) map[string]any {
	out := make(map[string]any)
	out["type"] = evt.Type()
	switch evt := evt.(type) {
	case *event.Tick:
		out["tick"] = evt.Tick

	case *event.StatBase:
		out["value"] = evt.Value
		out["setValue"] = func(value float64) {
			evt.Value = int(value)
		}
		statObj := ctx.GetObjectByID(evt.StatID)
		if statObj == nil {
			log.Println("WARNING: Error during EventToLua:", nbattle.ErrObjectNotFound)
			return out
		}
		stat, ok := statObj.(*nbattle.Stat)
		if !ok {
			log.Println("WARNING: Error during EventToLua:", nbattle.ErrUnexpectedObjectType)
			return out
		}
		out["stat"] = StatToLua(ctx, stat)

	case *event.StatMod:
		out["modValue"] = evt.ModValue
		out["setModValue"] = func(value float64) {
			evt.ModValue = int(value)
		}
		sourceCombatant, err := ctx.GetCombatantByID(evt.SourceID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		out["source"] = sourceCombatant

		ctx.GetObject(evt.StatID)

	case *event.NewCombatant:
		out["combatant"] = nil
		out["team"] = evt.Team
		combatant, err := ctx.GetCombatantByID(evt.ID)
		if err != nil {
			log.Println("WARNING: Error during EventToLua:", err)
			return out
		}
		out["combatant"] = CombatantToLua(ctx, combatant)

	}

	return out
}

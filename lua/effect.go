package lua

import (
	"io"
	"log"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/event"
	luago "github.com/rosbit/luago"
)

type LuaEffect struct {
	luaCtx *luago.LuaContext
}

func NewEffect(ctx *nbattle.Context, script io.Reader) (*nbattle.EffectDef, error) {
	scriptBytes, err := io.ReadAll(script)
	if err != nil {
		return nil, err
	}
	luaCtx, err := loadLuaScript(ctx, scriptBytes)
	if err != nil {
		return nil, err
	}
	nameIf, err := luaCtx.CallFunc("Name")
	if err != nil {
		return nil, err
	}
	name, ok := nameIf.(string)
	if !ok {
		return nil, ErrUnexpectedLuaFuncReturn
	}
	return ctx.NewEffectDef(name, func() nbattle.Effect {
		luaCtx, err := loadLuaScript(ctx, scriptBytes)
		if err != nil {
			return nil
		}
		return &LuaEffect{
			luaCtx: luaCtx,
		}
	}), nil
}

func loadLuaScript(ctx *nbattle.Context, scriptBytes []byte) (*luago.LuaContext, error) {
	luaCtx, err := luago.NewContext()
	if err != nil {
		return nil, err
	}
	if err := luaCtx.LoadScript(string(scriptBytes), luaGlobals(ctx)); err != nil {
		return nil, err
	}
	return luaCtx, nil
}

func (e *LuaEffect) OnAdd(ctx *nbattle.EffectCtx) {
	if _, err := e.luaCtx.CallFunc("OnAdd", effectCtxToLua(ctx)); err != nil {
		logLuaFuncCallError(err, ctx.Def.GetName()+".OnRemove")
	}
}

func (e *LuaEffect) OnRemove(ctx *nbattle.EffectCtx) {
	if _, err := e.luaCtx.CallFunc("OnRemove", effectCtxToLua(ctx)); err != nil {
		logLuaFuncCallError(err, ctx.Def.GetName()+".OnRemove")
	}
}

func (e *LuaEffect) OnEvent(ctx *nbattle.EffectCtx, evt event.Event) {
	switch evt := evt.(type) {
	case *event.Tick:
		if _, err := e.luaCtx.CallFunc("OnTick", effectCtxToLua(ctx), map[string]any{"tick": evt.Tick}); err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+".OnTick")
		}

	case *event.CombatantUpdate:
		combatant, err := ctx.Ctx.GetCombatantByID(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+".OnCombatantUpdate")
			break
		}
		if _, err := e.luaCtx.CallFunc("OnCombatantUpdate", effectCtxToLua(ctx), map[string]any{
			"combatant": combatantToLua(ctx.Ctx, combatant),
			"active":    evt.Active,
		}); err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+".OnCombatantUpdate")
		}

	case *event.CombatantStatBase:
		funcName := "OnCombatantStatBase"
		combatant, err := ctx.Ctx.GetCombatantByID(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		statDef, err := ctx.Ctx.GetStatDefByID(evt.StatDefID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"combatant": combatantToLua(ctx.Ctx, combatant),
			"statDef":   statDefToLua(statDef),
			"value":     evt.Value,
			"setValue": func(value int) {
				evt.Value = value
			},
		}); err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantStatMod:
		funcName := "OnCombatantStatMod"
		combatant, err := ctx.Ctx.GetCombatantByID(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		statDef, err := ctx.Ctx.GetStatDefByID(evt.StatDefID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"combatant": combatantToLua(ctx.Ctx, combatant),
			"statDef":   statDefToLua(statDef),
			"value":     evt.ModValue,
			"setValue": func(value int) {
				evt.ModValue = value
			},
		}); err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantEffect:
		funcName := "OnCombatantEffect"
		target, err := ctx.Ctx.GetCombatantByID(evt.TargetID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		effectDef, err := ctx.Ctx.GetEffectDefByID(evt.EffectDefID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		source, _ := ctx.Ctx.GetObject(evt.SourceID)
		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"target":  combatantToLua(ctx.Ctx, target),
			"effect":  effectDef.GetName(),
			"potency": evt.Potency,
			"source":  objectToLua(ctx.Ctx, source),
		}); err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
		}

	case *event.Trigger:
		funcName := "OnTrigger"

		triggerDefObj, err := ctx.Ctx.GetObject(evt.TriggerDefID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		triggerDef, ok := triggerDefObj.(*nbattle.TriggerDef)
		if !ok {
			logLuaFuncCallError(nbattle.ErrUnexpectedObjectType, ctx.Def.GetName()+"."+funcName)
			break
		}

		effectDef, err := ctx.Ctx.GetEffectDefByID(evt.EffectDefID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		effectTarget, err := ctx.Ctx.GetCombatantByID(evt.EffectTargetID)
		if err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
			break
		}
		effectSource, _ := ctx.Ctx.GetObject(evt.EffectSourceID)

		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"trigger": triggerDef.GetName(),
			"target":  combatantToLua(ctx.Ctx, effectTarget),
			"effect":  effectDef.GetName(),
			"potency": evt.EffectPotency,
			"source":  objectToLua(ctx.Ctx, effectSource),
		}); err != nil {
			logLuaFuncCallError(err, ctx.Def.GetName()+"."+funcName)
		}

	default:
		log.Printf("WARNING: Effect %s received unknown event type: %T", ctx.Def.GetName(), evt.Type())
	}
}

func effectCtxToLua(effectCtx *nbattle.EffectCtx) map[string]any {
	return map[string]any{
		"target":      combatantToLua(effectCtx.Ctx, effectCtx.Target),
		"source":      objectToLua(effectCtx.Ctx, effectCtx.Source),
		"effect":      effectCtx.Def.GetName(),
		"potency":     effectCtx.Potency,
		"emitTrigger": effectCtx.EmitTrigger,
	}
}

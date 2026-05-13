package nbattle

import (
	"fmt"
	"io"
	"strings"

	"github.com/chompy/nbattle_go/event"
	luago "github.com/rosbit/luago"
)

type luaEffect struct {
	luaCtx *luago.LuaContext
}

// NewEffect creates a new EffectDef from a Lua script.
func (c *Context) NewLuaEffect(script io.Reader) (*EffectDef, error) {
	scriptBytes, err := io.ReadAll(script)
	if err != nil {
		return nil, err
	}
	luaCtx, err := loadLuaScript(c, scriptBytes)
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
	return c.NewEffectDef(name, func() Effect {
		luaCtx, err := loadLuaScript(c, scriptBytes)
		if err != nil {
			return nil
		}
		return &luaEffect{
			luaCtx: luaCtx,
		}
	}), nil
}

func luaGlobals(ctx *Context) map[string]any {

	globals := make(map[string]any)
	globals["ctx"] = contextToLua(ctx)
	for _, statDef := range ctx.GetStatDefs() {
		globals["STAT_"+strings.ToUpper(statDef.GetName())] = statDefToLua(ctx, statDef)
	}
	for name, value := range ctx.GetFlags() {
		globals["FLAG_"+strings.ToUpper(name)] = value
	}
	return globals
}

func contextToLua(ctx *Context) map[string]any {
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

func errorToLua(ctx *Context, err error) map[string]any {
	ctx.log.Warn("Error during Lua call.", "error", err)
	return map[string]any{
		"type":  ObjectTypeError,
		"error": err.Error(),
	}
}

func objectIDFromLua(object any) (int, error) {
	switch object := object.(type) {
	case map[string]any:
		id, ok := object["id"].(int)
		if !ok {
			return 0, ErrUnexpectedObjectType
		}
		return id, nil

	case float64:
		return int(object), nil

	case float32:
		return int(object), nil

	case int:
		return object, nil

	}
	return 0, ErrUnexpectedObjectType
}

func objectFromLua(ctx *Context, object any) (Object, error) {
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

func objectToLua(ctx *Context, object any) map[string]any {
	switch object := object.(type) {
	case *StatDef:
		return statDefToLua(ctx, object)

	case *EffectDef:
		return map[string]any{
			"id":   object.GetID(),
			"type": object.GetType(),
			"name": object.GetName(),
		}

	case *TriggerDef:
		return map[string]any{
			"id":   object.GetID(),
			"type": object.GetType(),
			"name": object.GetName(),
		}

	case *Stat:
		return statToLua(ctx, object)

	case *Combatant:
		return combatantToLua(ctx, object)

	default:
		return errorToLua(ctx, ErrUnexpectedObjectType)
	}
}

func combatantToLua(ctx *Context, combatant *Combatant) map[string]any {
	if combatant == nil {
		return map[string]any{
			"type":  ObjectTypeUnknown,
			"error": ErrNilObject,
		}
	}
	return map[string]any{
		"id":   combatant.GetID(),
		"type": combatant.GetType(),
		"getStat": func(statDef any) map[string]any {
			stat, err := combatant.GetStat(statDef)
			if err != nil {
				logLuaFuncCallError(ctx.log, err, fmt.Sprintf("combatant.%d.getStat", combatant.GetID()))
				return errorToLua(ctx, err)
			}
			return statToLua(ctx, stat)
		},
		"setEffect": func(effectDef any, potency int, sourceObj any) {
			if err := combatant.SetEffect(effectDef, potency, sourceObj); err != nil {
				logLuaFuncCallError(ctx.log, err, fmt.Sprintf("combatant.%d.setEffect", combatant.GetID()))
			}
		},
		"hasEffect": combatant.HasEffect,
		"removeEffect": func(effectDef any) {
			if err := combatant.SetEffect(effectDef, 0, nil); err != nil {
				logLuaFuncCallError(ctx.log, err, fmt.Sprintf("combatant.%d.removeEffect", combatant.GetID()))
			}
		},
		"setFlag": combatant.SetFlag,
		"hasFlag": combatant.HasFlag,
	}
}

func loadLuaScript(ctx *Context, scriptBytes []byte) (*luago.LuaContext, error) {
	luaCtx, err := luago.NewContext()
	if err != nil {
		return nil, err
	}
	if err := luaCtx.LoadScript(string(scriptBytes), luaGlobals(ctx)); err != nil {
		return nil, err
	}
	return luaCtx, nil
}

func (e *luaEffect) OnAdd(ctx *EffectCtx) {
	if _, err := e.luaCtx.CallFunc("OnAdd", effectCtxToLua(ctx)); err != nil {
		logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+".OnAdd")
	}
}

func (e *luaEffect) OnRemove(ctx *EffectCtx) {
	if _, err := e.luaCtx.CallFunc("OnRemove", effectCtxToLua(ctx)); err != nil {
		logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+".OnRemove")
	}
}

func (e *luaEffect) OnEvent(ctx *EffectCtx, evt event.Event) {
	switch evt := evt.(type) {
	case *event.Tick:
		if _, err := e.luaCtx.CallFunc("OnTick", effectCtxToLua(ctx), map[string]any{"tick": evt.Tick}); err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+".OnTick")
		}

	case *event.CombatantUpdate:
		combatant, err := ctx.Ctx.GetCombatant(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+".OnCombatantUpdate")
			break
		}
		if _, err := e.luaCtx.CallFunc("OnCombatantUpdate", effectCtxToLua(ctx), map[string]any{
			"combatant": combatantToLua(ctx.Ctx, combatant),
			"active":    evt.Active,
			"flags":     evt.Flags,
		}); err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+".OnCombatantUpdate")
		}

	case *event.CombatantStatBase:
		funcName := "OnCombatantStatBase"
		combatant, err := ctx.Ctx.GetCombatant(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		statDef, err := ctx.Ctx.GetStatDef(evt.StatDefID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"combatant": combatantToLua(ctx.Ctx, combatant),
			"statDef":   statDefToLua(ctx.Ctx, statDef),
			"value":     evt.Value,
			"setValue": func(value int) {
				evt.Value = value
			},
		}); err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantStatMod:
		funcName := "OnCombatantStatMod"
		combatant, err := ctx.Ctx.GetCombatant(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		statDef, err := ctx.Ctx.GetStatDef(evt.StatDefID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"combatant": combatantToLua(ctx.Ctx, combatant),
			"statDef":   statDefToLua(ctx.Ctx, statDef),
			"value":     evt.ModValue,
			"setValue": func(value int) {
				evt.ModValue = value
			},
		}); err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantEffect:
		funcName := "OnCombatantEffect"
		target, err := ctx.Ctx.GetCombatant(evt.TargetID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		effectDef, err := ctx.Ctx.GetEffectDef(evt.EffectDefID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		source, _ := ctx.Ctx.GetObject(evt.SourceID)
		if _, err := e.luaCtx.CallFunc(funcName, effectCtxToLua(ctx), map[string]any{
			"target":  combatantToLua(ctx.Ctx, target),
			"effect":  effectDef.GetName(),
			"potency": evt.Potency,
			"source":  objectToLua(ctx.Ctx, source),
		}); err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
		}

	case *event.Trigger:
		funcName := "OnTrigger"
		triggerDef, err := ctx.Ctx.GetTriggerDef(evt.TriggerDefID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		effectDef, err := ctx.Ctx.GetEffectDef(evt.EffectDefID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
			break
		}
		effectTarget, err := ctx.Ctx.GetCombatant(evt.EffectTargetID)
		if err != nil {
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
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
			logLuaFuncCallError(ctx.Ctx.log, err, ctx.Def.GetName()+"."+funcName)
		}

	default:
		ctx.Ctx.log.Warn("Effect received unknown event type.", "effect", ctx.Def, "event", evt.Type())
	}
}

func effectCtxToLua(effectCtx *EffectCtx) map[string]any {
	if effectCtx == nil || effectCtx.Ctx == nil {
		return map[string]any{
			"type":  ObjectTypeError,
			"error": ErrNilObject,
		}
	}
	return map[string]any{
		"target":      combatantToLua(effectCtx.Ctx, effectCtx.Target),
		"source":      objectToLua(effectCtx.Ctx, effectCtx.Source),
		"effect":      effectCtx.Def.GetName(),
		"potency":     effectCtx.Potency,
		"emitTrigger": effectCtx.EmitTrigger,
		"remove": func() {
			effectCtx.Target.SetEffect(effectCtx.Def, 0, nil)
		},
	}
}

func statDefToLua(ctx *Context, statDef *StatDef) map[string]any {
	if statDef == nil {
		return errorToLua(ctx, ErrNilObject)
	}
	return map[string]any{
		"type": statDef.GetType(),
		"name": statDef.GetName(),
		"min":  statDef.GetMin(),
		"max":  statDef.GetMax(),
	}
}

func statToLua(ctx *Context, stat *Stat) map[string]any {
	if stat == nil {
		return errorToLua(ctx, ErrNilObject)
	}
	return map[string]any{
		"def":     statDefToLua(ctx, stat.GetDef()),
		"getBase": stat.GetBase,
		"setBase": func(value float64) {
			stat.SetBase(int(value))
		},
		"addBase": func(value float64) {
			stat.AddBase(int(value))
		},
		"subBase": func(value float64) {
			stat.SubBase(int(value))
		},
		"getValue": stat.GetValue,
		"setMod": func(source any, value int) {
			sourceObj, err := objectFromLua(ctx, source)
			if err != nil {
				logLuaFuncCallError(ctx.log, err, "stat.setMod")
				return
			}
			stat.SetMod(sourceObj, value)
		},
		"getCombatant": func() map[string]any {
			combatant, err := ctx.GetCombatantWithStat(stat)
			if err != nil {
				return errorToLua(ctx, err)
			}
			return combatantToLua(ctx, combatant)
		},
	}
}

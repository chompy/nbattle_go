package nbattle

import (
	"fmt"
	"io"
	"strings"

	"github.com/chompy/nbattle_go/event"
	luago "github.com/rosbit/luago"
)

type luaDeferedEffect struct {
	def     *EffectDef
	target  *Combatant
	source  Object
	potency int
}

type luaEffect struct {
	ctx            *Context
	luaCtx         *luago.LuaContext
	deferedEffects []*luaDeferedEffect
}

// NewEffect creates a new EffectDef from a Lua script.
func (c *Context) NewLuaEffect(script io.Reader) (*EffectDef, error) {
	scriptBytes, err := io.ReadAll(script)
	if err != nil {
		return nil, err
	}
	luaCtx, err := luago.NewContext()
	if err != nil {
		return nil, err
	}
	if err := luaCtx.LoadScript(string(scriptBytes), nil); err != nil {
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
	return c.NewEffectDef(strings.ToLower(strings.Clone(name)), func() Effect {
		newLuaEffect := &luaEffect{c, nil, make([]*luaDeferedEffect, 0)}
		luaCtx, err := luago.NewContext()
		if err != nil {
			return nil
		}
		if err := luaCtx.LoadScript(string(scriptBytes), newLuaEffect.luaGlobals()); err != nil {
			return nil
		}
		newLuaEffect.luaCtx = luaCtx
		return newLuaEffect
	}), nil
}

func (l *luaEffect) luaGlobals() map[string]any {
	globals := make(map[string]any)
	for _, statDef := range l.ctx.GetStatDefs() {
		globals["STAT_"+strings.ToUpper(statDef.GetName())] = l.statDefToLua(statDef)
	}
	for name, value := range l.ctx.GetFlags() {
		globals["FLAG_"+strings.ToUpper(name)] = value
	}
	globals["ctx"] = map[string]any{
		"getTick": l.ctx.GetTick,
		"getCombatants": func() []map[string]any {
			out := make([]map[string]any, 0)
			for _, combatant := range l.ctx.GetCombatants() {
				out = append(out, l.combatantToLua(combatant))
			}
			return out
		},
		"getObject": func(obj any) map[string]any {
			return l.objectToLua(obj)
		},
		"emitTrigger": func(trigger any, source any) {
			l.ctx.EmitTrigger(trigger, source)
		},
	}

	return globals
}

func (l *luaEffect) errorToLua(err error) map[string]any {
	l.ctx.log.Warn("Error during Lua call.", "error", err)
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

func (l *luaEffect) objectFromLua(object any) (Object, error) {
	id, err := objectIDFromLua(object)
	if err != nil {
		return nil, err
	}
	obj, err := l.ctx.GetObjectByID(id)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (l *luaEffect) objectToLua(object any) map[string]any {
	switch object := object.(type) {
	case *StatDef:
		return l.statDefToLua(object)

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
		return l.statToLua(object)

	case *Combatant:
		return l.combatantToLua(object)

	default:
		return l.errorToLua(ErrUnexpectedObjectType)
	}
}

func (l *luaEffect) statDefToLua(statDef *StatDef) map[string]any {
	if statDef == nil {
		return l.errorToLua(ErrNilObject)
	}
	return map[string]any{
		"type": statDef.GetType(),
		"name": statDef.GetName(),
		"min":  statDef.GetMin(),
		"max":  statDef.GetMax(),
	}
}

func (l *luaEffect) statToLua(stat *Stat) map[string]any {
	if stat == nil {
		return l.errorToLua(ErrNilObject)
	}
	return map[string]any{
		"def":     l.statDefToLua(stat.GetDef()),
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
			sourceObj, err := l.objectFromLua(source)
			if err != nil {
				logLuaFuncCallError(l.ctx.log, err, "stat.setMod")
				return
			}
			stat.SetMod(sourceObj, value)
		},
		"getCombatant": func() map[string]any {
			combatant, err := l.ctx.GetCombatantWithStat(stat)
			if err != nil {
				return l.errorToLua(err)
			}
			return l.combatantToLua(combatant)
		},
	}
}

func (l *luaEffect) combatantToLua(combatant *Combatant) map[string]any {
	if combatant == nil {
		return map[string]any{
			"type":  ObjectTypeUnknown,
			"error": ErrNilObject,
		}
	}
	setEffect := func(effectDefObj any, sourceObj any, potency int) map[string]any {
		effectDef, err := l.ctx.GetEffectDef(effectDefObj)
		if err != nil {
			return l.errorToLua(err)
		}
		source, err := l.ctx.GetObject(sourceObj)
		if err != nil {
			return l.errorToLua(err)
		}
		l.deferedEffects = append(l.deferedEffects, &luaDeferedEffect{effectDef, combatant, source, potency})
		return nil
	}
	return map[string]any{
		"id":   combatant.GetID(),
		"type": combatant.GetType(),
		"getStat": func(statDef any) map[string]any {
			stat, err := combatant.GetStat(statDef)
			if err != nil {
				logLuaFuncCallError(l.ctx.log, err, fmt.Sprintf("combatant.%d.getStat", combatant.GetID()))
				return l.errorToLua(err)
			}
			return l.statToLua(stat)
		},
		"setEffect": setEffect,
		"hasEffect": combatant.HasEffect,
		"removeEffect": func(effectDefObj any, sourceObj any) map[string]any {
			return setEffect(effectDefObj, sourceObj, 0)
		},
		"setFlag": combatant.SetFlag,
		"hasFlag": combatant.HasFlag,
	}
}

func (l *luaEffect) effectCtxToLua(effectCtx *EffectContext) map[string]any {
	if effectCtx == nil {
		return map[string]any{
			"type":  ObjectTypeError,
			"error": ErrNilObject,
		}
	}
	return map[string]any{
		"name":    effectCtx.Def.GetName(),
		"target":  l.combatantToLua(effectCtx.Target),
		"source":  l.objectToLua(effectCtx.Source),
		"potency": effectCtx.Potency,
		"remove": func() {
			l.deferedEffects = append(l.deferedEffects, &luaDeferedEffect{effectCtx.Def, effectCtx.Target, effectCtx.Source, 0})
		},
	}
}

func (l *luaEffect) processDeferedEffects() error {
	effects := make([]*luaDeferedEffect, len(l.deferedEffects))
	copy(effects, l.deferedEffects)
	l.deferedEffects = make([]*luaDeferedEffect, 0)
	for _, effect := range effects {
		if err := effect.target.SetEffect(effect.def, effect.source, effect.potency); err != nil {
			return err
		}
	}
	return nil
}

func (l *luaEffect) OnAdd(ctx *Context, effectCtx *EffectContext) {
	if _, err := l.luaCtx.CallFunc("OnAdd", l.effectCtxToLua(effectCtx)); err != nil {
		logLuaFuncCallError(ctx.log, err, effectCtx.Def.GetName()+".OnAdd")
	}
	if err := l.processDeferedEffects(); err != nil {
		logLuaFuncCallError(ctx.log, err, effectCtx.Def.GetName()+".OnAdd")
	}
}

func (l *luaEffect) OnRemove(ctx *Context, effectCtx *EffectContext) {
	if _, err := l.luaCtx.CallFunc("OnRemove", l.effectCtxToLua(effectCtx)); err != nil {
		logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+".OnRemove")
	}
	if err := l.processDeferedEffects(); err != nil {
		logLuaFuncCallError(ctx.log, err, effectCtx.Def.GetName()+".OnAdd")
	}
}

func (l *luaEffect) OnEvent(ctx *Context, effectCtx *EffectContext, evt event.Event) {
	switch evt := evt.(type) {
	case *event.Tick:
		if _, err := l.luaCtx.CallFunc("OnTick", l.effectCtxToLua(effectCtx), map[string]any{"tick": evt.Tick}); err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+".OnTick")
		}

	case *event.CombatantUpdate:
		funcName := "OnCombatantUpdate"
		combatant, err := l.ctx.GetCombatant(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := l.luaCtx.CallFunc(funcName, l.effectCtxToLua(effectCtx), map[string]any{
			"combatant": l.combatantToLua(combatant),
			"active":    evt.Active,
			"flags":     evt.Flags,
		}); err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantStatBase:
		funcName := "OnCombatantStatBase"
		combatant, err := l.ctx.GetCombatant(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		statDef, err := l.ctx.GetStatDef(evt.StatDefID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := l.luaCtx.CallFunc(funcName, l.effectCtxToLua(effectCtx), map[string]any{
			"combatant": l.combatantToLua(combatant),
			"statDef":   l.statDefToLua(statDef),
			"value":     evt.Value,
			"setValue": func(value int) {
				evt.Value = value
			},
		}); err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantStatMod:
		funcName := "OnCombatantStatMod"
		combatant, err := l.ctx.GetCombatant(evt.CombatantID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		statDef, err := l.ctx.GetStatDef(evt.StatDefID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := l.luaCtx.CallFunc(funcName, l.effectCtxToLua(effectCtx), map[string]any{
			"combatant": l.combatantToLua(combatant),
			"statDef":   l.statDefToLua(statDef),
			"value":     evt.ModValue,
			"setValue": func(value int) {
				evt.ModValue = value
			},
		}); err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
		}

	case *event.CombatantEffect:
		funcName := "OnCombatantEffect"
		target, err := l.ctx.GetCombatant(evt.TargetID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		effectDef, err := l.ctx.GetEffectDef(evt.EffectDefID)
		if err != nil {
			logLuaFuncCallError(ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		source, _ := ctx.GetObject(evt.SourceID)
		if _, err := l.luaCtx.CallFunc(funcName, l.effectCtxToLua(effectCtx), map[string]any{
			"target":  l.combatantToLua(target),
			"effect":  effectDef.GetName(),
			"potency": evt.Potency,
			"source":  l.objectToLua(source),
		}); err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
		}

	case *event.Trigger:
		funcName := "OnTrigger"
		triggerDef, err := l.ctx.GetTriggerDef(evt.TriggerDefID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		source, err := l.ctx.GetObject(evt.SourceID)
		if err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
			break
		}
		if _, err := l.luaCtx.CallFunc(funcName, l.effectCtxToLua(effectCtx), map[string]any{
			"trigger": triggerDef.GetName(),
			"source":  l.objectToLua(source),
		}); err != nil {
			logLuaFuncCallError(l.ctx.log, err, effectCtx.Def.GetName()+"."+funcName)
		}

	default:
		ctx.log.Warn("Effect received unknown event type.", "effect", effectCtx.Def, "event", evt.Type())
	}

	if err := l.processDeferedEffects(); err != nil {
		logLuaFuncCallError(ctx.log, err, effectCtx.Def.GetName()+".OnAdd")
	}

}

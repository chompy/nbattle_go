package lua

import (
	"io"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/internal/event"
	luago "github.com/rosbit/luago"
)

type LuaEffect struct {
	luaCtx *luago.LuaContext
}

func loadLuaScript(ctx *nbattle.Context, scriptBytes []byte) (*luago.LuaContext, error) {
	luaCtx, err := luago.NewContext()
	if err != nil {
		return nil, err
	}
	if err := luaCtx.LoadScript(string(scriptBytes), LuaGlobals(ctx)); err != nil {
		return nil, err
	}
	return luaCtx, nil
}

func NewLuaEffect(ctx *nbattle.Context, script io.Reader) (*nbattle.EffectDef, error) {
	scriptBytes, err := io.ReadAll(script)
	if err != nil {
		return nil, err
	}
	luaCtx, err := loadLuaScript(ctx, scriptBytes)
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

func (e *LuaEffect) OnAdd(ctx *nbattle.EffectCtx) {
	if _, err := e.luaCtx.CallFunc("OnAdd",
		CombatantToLua(ctx.Ctx, ctx.Target),
		CombatantToLua(ctx.Ctx, ctx.Source),
	); err != nil {
		logLuaFuncCallError(err, ctx.Def.GetName()+".OnAdd")
	}
}

func (e *LuaEffect) OnRemove(ctx *nbattle.EffectCtx) {
	if _, err := e.luaCtx.CallFunc("OnRemove",
		CombatantToLua(ctx.Ctx, ctx.Target),
		CombatantToLua(ctx.Ctx, ctx.Source),
	); err != nil {
		logLuaFuncCallError(err, ctx.Def.GetName()+".OnRemove")
	}
}

func (e *LuaEffect) OnEvent(ctx *nbattle.EffectCtx, event event.Event) {
	if _, err := e.luaCtx.CallFunc("OnEvent",
		EventToLua(ctx.Ctx, event),
		CombatantToLua(ctx.Ctx, ctx.Target),
		CombatantToLua(ctx.Ctx, ctx.Source),
	); err != nil {
		logLuaFuncCallError(err, ctx.Def.GetName()+".OnEvent")
	}
}

func EffectContextToLua(ctx *nbattle.EffectCtx) map[string]any {
	return map[string]any{
		"source": CombatantToLua(ctx.Ctx, ctx.Source),
		"target": CombatantToLua(ctx.Ctx, ctx.Target),
	}
}

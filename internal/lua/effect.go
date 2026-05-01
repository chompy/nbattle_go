package lua

import (
	nbattle "github.com/chompy/nbattle_go/nbattle_go"
	luago "github.com/rosbit/luago"
)

type LuaEffect struct {
	luaCtx *luago.LuaContext
}

func (e *LuaEffect) OnAdd(ctx *nbattle.EffectCtx) {

	e.luaCtx.CallFunc("OnAdd", map[string]any{
		"source": ctx.Source.GetID(),
		"target": ctx.Target.GetID(),
	})

}

func (e *LuaEffect) OnRemove(ctx *nbattle.EffectCtx) {

}

func (e *LuaEffect) OnEvent(ctx *nbattle.EffectCtx, event *nbattle.Event) {

}

/*
	OnAdd(ctx *EffectCtx)
	OnRemove(ctx *EffectCtx)
	OnEvent(ctx *EffectCtx, event *Event)
*/

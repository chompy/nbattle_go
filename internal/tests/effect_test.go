package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/internal/event"
)

type TestEffect struct {
	hpStat *nbattle.StatDef
}

func (e *TestEffect) OnAdd(ctx *nbattle.EffectCtx) {
	ctx.Target.Stat(e.hpStat).SetBase(29)
}

func (e *TestEffect) OnRemove(ctx *nbattle.EffectCtx) {
	ctx.Target.Stat(e.hpStat).SetBase(31)
}

func (e *TestEffect) OnEvent(ctx *nbattle.EffectCtx, evt event.Event) {
	if evt.Type() == event.TickEvent {
		ctx.Target.Stat(e.hpStat).AddBase(-1)
	}
}

func TestCombatantEffect(t *testing.T) {

	ctx := nbattle.New()

	statDefHP := ctx.NewStatDef("hp", 0, 99)
	effectDefTest := ctx.NewEffectDef("test", func() nbattle.Effect {
		return &TestEffect{hpStat: statDefHP}
	})

	cmbt := ctx.NewCombatant(1)
	cmbt.Stat(statDefHP).SetBase(15)

	srcCmbt := ctx.NewCombatant(1)

	cmbt.AddEffect(effectDefTest, srcCmbt)

	if cmbt.Stat(statDefHP).GetValue() != 29 {
		t.Fatal("expected effect to set combatant hp to 29 on add")
	}

	ctx.Tick()

	if cmbt.Stat(statDefHP).GetValue() != 28 {
		t.Fatal("expected effect reduce combatant hp by 1 on tick")
	}

	ctx.Tick()

	if cmbt.Stat(statDefHP).GetValue() != 27 {
		t.Fatal("expected effect reduce combatant hp by 1 on tick")
	}

	cmbt.RemoveEffect(effectDefTest)

	if cmbt.Stat(statDefHP).GetValue() != 31 {
		t.Fatal("expected effect to set combatant hp to 31 on removal")
	}

}

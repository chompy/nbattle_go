package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go/nbattle_go"
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

func (e *TestEffect) OnEvent(ctx *nbattle.EffectCtx, event *nbattle.Event) {
	if event.Type() == nbattle.EventTypeTick {
		ctx.Target.Stat(e.hpStat).AddBase(-1)
	}
}

func TestCombatantEffect(t *testing.T) {

	ctx := nbattle.New()

	statDefHP := ctx.NewStatDef(0, 99)
	effectDefTest := ctx.NewEffectDef(func() nbattle.Effect {
		return &TestEffect{hpStat: statDefHP}
	})

	cmbt := ctx.NewCombatant()
	cmbt.Stat(statDefHP).SetBase(15)

	srcCmbt := ctx.NewCombatant()

	cmbt.AddEffect(effectDefTest, srcCmbt)

	if cmbt.Stat(statDefHP).Value() != 29 {
		t.Fatal("expected effect to set combatant hp to 29 on add")
	}

	ctx.Tick()

	if cmbt.Stat(statDefHP).Value() != 28 {
		t.Fatal("expected effect reduce combatant hp by 1 on tick")
	}

	ctx.Tick()

	if cmbt.Stat(statDefHP).Value() != 27 {
		t.Fatal("expected effect reduce combatant hp by 1 on tick")
	}

	cmbt.RemoveEffect(effectDefTest)

	if cmbt.Stat(statDefHP).Value() != 31 {
		t.Fatal("expected effect to set combatant hp to 31 on removal")
	}

}

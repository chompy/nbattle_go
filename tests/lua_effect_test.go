package tests

import (
	"embed"
	"log"
	"testing"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/event"
)

//go:embed effects/*.lua
var luaEffectFS embed.FS

func TestPoisonEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/poison.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(50)

	attacker := ctx.NewCombatant()
	target.SetEffect(effectDef, attacker, 1)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 50 {
		t.Fatal("expected hp to be 50 before poison tick")
	}

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 48 {
		t.Fatalf("expected hp to be 48 after 1 poison tick, got %d", hp.GetValue())
	}

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 46 {
		t.Fatalf("expected hp to be 46 after 2 poison ticks, got %d", hp.GetValue())
	}

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 44 {
		t.Fatalf("expected hp to be 44 after 3 poison ticks, got %d", hp.GetValue())
	}
}

func TestRegenerateEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/regenerate.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(10)

	target.SetEffect(effectDef, target, 1)

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 13 {
		t.Fatalf("expected hp to be 13 after 1 regen tick, got %d", hp.GetValue())
	}

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 16 {
		t.Fatalf("expected hp to be 16 after 2 regen ticks, got %d", hp.GetValue())
	}

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 19 {
		t.Fatalf("expected hp to be 19 after 3 regen ticks, got %d", hp.GetValue())
	}
}

func TestPoisonAndRegenerateInteraction(t *testing.T) {
	poisonF, err := luaEffectFS.Open("effects/poison.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer poisonF.Close()

	regenF, err := luaEffectFS.Open("effects/regenerate.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer regenF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 200)

	poisonDef, err := ctx.NewLuaEffect(poisonF)
	if err != nil {
		t.Fatal(err)
	}

	regenDef, err := ctx.NewLuaEffect(regenF)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(100)

	attacker := ctx.NewCombatant()
	target.SetEffect(poisonDef, attacker, 1)
	target.SetEffect(regenDef, target, 1)

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 101 {
		t.Fatalf("expected hp to be 101 after 1 tick (100 - 2 poison + 3 regen = 101), got %d", hp.GetValue())
	}

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 102 {
		t.Fatalf("expected hp to be 102 after 2 ticks, got %d", hp.GetValue())
	}

	target.SetEffect(poisonDef, attacker, 0)

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 105 {
		t.Fatalf("expected hp to be 105 after removing poison (only regen), got %d", hp.GetValue())
	}
}

func TestShieldEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/shield.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 200)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	target.SetEffect(effectDef, target, 1)
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(10)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 20 {
		t.Fatalf("expected hp to be 20 (doubled by shield), got %d", hp.GetValue())
	}

	hp, _ = target.GetStat(hpStatDef)
	hp.SetBase(5)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 10 {
		t.Fatalf("expected hp to be 10 (5*2), got %d", hp.GetValue())
	}

	hp, _ = target.GetStat(hpStatDef)
	hp.SetBase(15)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 30 {
		t.Fatalf("expected hp to be 30 (15*2), got %d", hp.GetValue())
	}
}

func TestCounterEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/counter.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 200)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(30)

	attacker := ctx.NewCombatant()
	attHp, _ := attacker.GetStat(hpStatDef)
	attHp.SetBase(30)

	target.SetEffect(effectDef, attacker, 1)

	targetHp, _ := target.GetStat(hpStatDef)
	targetHp.SetBase(25)

	targetHp, _ = target.GetStat(hpStatDef)
	if targetHp.GetValue() != 25 {
		t.Fatalf("expected target hp to be 25, got %d", targetHp.GetValue())
	}

	attHp, _ = attacker.GetStat(hpStatDef)
	if attHp.GetValue() != 25 {
		t.Fatalf("expected attacker hp to be 25 (30-5 counter), got %d", attHp.GetValue())
	}
}

func TestSelfHealEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/self_heal.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(20)

	target.SetEffect(effectDef, target, 1)

	hp, _ = target.GetStat(hpStatDef)
	hp.SetBase(15)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 16 {
		t.Fatalf("expected hp to be 16 (15 + 1 self heal), got %d", hp.GetValue())
	}

	hp, _ = target.GetStat(hpStatDef)
	hp.SetBase(10)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 11 {
		t.Fatalf("expected hp to be 11 (10 + 1 self heal), got %d", hp.GetValue())
	}

	hp, _ = target.GetStat(hpStatDef)
	hp.SetBase(20)

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 20 {
		t.Fatalf("expected hp to be 20 (full, no heal), got %d", hp.GetValue())
	}
}

func TestTriggerEffect(t *testing.T) {
	triggerF, err := luaEffectFS.Open("effects/trigger_effect.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer triggerF.Close()

	buffF, err := luaEffectFS.Open("effects/buff.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer buffF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)
	strStatDef := ctx.NewStatDef("str", 0, 99)

	triggerDef, err := ctx.NewLuaEffect(triggerF)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.NewLuaEffect(buffF)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(30)
	str, _ := target.GetStat(strStatDef)
	str.SetBase(5)

	attacker := ctx.NewCombatant()
	target.SetEffect(triggerDef, attacker, 1)

	hp, _ = target.GetStat(hpStatDef)
	hp.SetBase(0)

	str, _ = target.GetStat(strStatDef)
	if str.GetValue() != 15 {
		t.Fatalf("expected str to be 15 (5 + 10 buff), got %d", str.GetValue())
	}
}

func TestCopyStatEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/copy_stat.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)
	strStatDef := ctx.NewStatDef("str", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	source := ctx.NewCombatant()
	srcHp, _ := source.GetStat(hpStatDef)
	srcHp.SetBase(30)
	srcStr, _ := source.GetStat(strStatDef)
	srcStr.SetBase(12)

	target := ctx.NewCombatant()
	tgtHp, _ := target.GetStat(hpStatDef)
	tgtHp.SetBase(20)
	tgtStr, _ := target.GetStat(strStatDef)
	tgtStr.SetBase(3)

	target.SetEffect(effectDef, source, 1)

	tgtStr, _ = target.GetStat(strStatDef)
	if tgtStr.GetValue() != 12 {
		t.Fatalf("expected str to be 12 (copied from source), got %d", tgtStr.GetValue())
	}

	tgtHp, _ = target.GetStat(hpStatDef)
	if tgtHp.GetValue() != 20 {
		t.Fatalf("expected hp to still be 20 (unchanged), got %d", tgtHp.GetValue())
	}
}

func TestReflectEffect(t *testing.T) {
	f, err := luaEffectFS.Open("effects/reflect.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 200)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	tgtHp, _ := target.GetStat(hpStatDef)
	tgtHp.SetBase(30)

	attacker := ctx.NewCombatant()
	attHp, _ := attacker.GetStat(hpStatDef)
	attHp.SetBase(30)

	target.SetEffect(effectDef, attacker, 1)

	tgtHp, _ = target.GetStat(hpStatDef)
	tgtHp.SetBase(20)

	tgtHp, _ = target.GetStat(hpStatDef)
	if tgtHp.GetValue() != 20 {
		t.Fatalf("expected target hp to be 20, got %d", tgtHp.GetValue())
	}

	attHp, _ = attacker.GetStat(hpStatDef)
	if attHp.GetValue() != 20 {
		t.Fatalf("expected attacker hp to be 20 (30-10 reflected), got %d", attHp.GetValue())
	}
}

func TestMultipleCombatantsWithEffects(t *testing.T) {
	poisonF, err := luaEffectFS.Open("effects/poison.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer poisonF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	poisonDef, err := ctx.NewLuaEffect(poisonF)
	if err != nil {
		t.Fatal(err)
	}

	cmbt1 := ctx.NewCombatant()
	hp1, _ := cmbt1.GetStat(hpStatDef)
	hp1.SetBase(50)

	cmbt2 := ctx.NewCombatant()
	hp2, _ := cmbt2.GetStat(hpStatDef)
	hp2.SetBase(40)

	cmbt3 := ctx.NewCombatant()
	hp3, _ := cmbt3.GetStat(hpStatDef)
	hp3.SetBase(60)

	cmbt1.SetEffect(poisonDef, cmbt2, 1)
	cmbt2.SetEffect(poisonDef, cmbt3, 1)

	ctx.Tick()

	hp1, _ = cmbt1.GetStat(hpStatDef)
	if hp1.GetValue() != 48 {
		t.Fatalf("expected cmbt1 hp to be 48, got %d", hp1.GetValue())
	}

	hp2, _ = cmbt2.GetStat(hpStatDef)
	if hp2.GetValue() != 38 {
		t.Fatalf("expected cmbt2 hp to be 38, got %d", hp2.GetValue())
	}

	hp3, _ = cmbt3.GetStat(hpStatDef)
	if hp3.GetValue() != 60 {
		t.Fatalf("expected cmbt3 hp to be 60 (no poison), got %d", hp3.GetValue())
	}

	ctx.Tick()

	hp1, _ = cmbt1.GetStat(hpStatDef)
	if hp1.GetValue() != 46 {
		t.Fatalf("expected cmbt1 hp to be 46, got %d", hp1.GetValue())
	}

	hp2, _ = cmbt2.GetStat(hpStatDef)
	if hp2.GetValue() != 36 {
		t.Fatalf("expected cmbt2 hp to be 36, got %d", hp2.GetValue())
	}
}

func TestNewCombatantEvent(t *testing.T) {
	f, err := luaEffectFS.Open("effects/copy_stat.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)
	strStatDef := ctx.NewStatDef("str", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	source := ctx.NewCombatant()
	srcHp, _ := source.GetStat(hpStatDef)
	srcHp.SetBase(30)
	srcStr, _ := source.GetStat(strStatDef)
	srcStr.SetBase(15)

	target := ctx.NewCombatant()
	target.SetEffect(effectDef, source, 1)

	tgtStr, _ := target.GetStat(strStatDef)
	if tgtStr.GetValue() != 15 {
		t.Fatalf("expected str to be 15 (copied from source), got %d", tgtStr.GetValue())
	}
}

func TestEffectOnNilSource(t *testing.T) {
	f, err := luaEffectFS.Open("effects/regenerate.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(10)

	err = target.SetEffect(effectDef, nil, 1)
	if err == nil {
		t.Fatal("expected error when adding effect with nil source")
	}
}

func TestEffectRemoval(t *testing.T) {
	regenF, err := luaEffectFS.Open("effects/regenerate.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer regenF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	effectDef, err := ctx.NewLuaEffect(regenF)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(10)

	target.SetEffect(effectDef, target, 1)

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 13 {
		t.Fatalf("expected hp to be 13 after regen, got %d", hp.GetValue())
	}

	target.SetEffect(effectDef, target, 0)

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 13 {
		t.Fatalf("expected hp to stay 13 after removing regen, got %d", hp.GetValue())
	}
}

func TestMultipleEffectsSameCombatant(t *testing.T) {
	poisonF, err := luaEffectFS.Open("effects/poison.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer poisonF.Close()

	regenF, err := luaEffectFS.Open("effects/regenerate.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer regenF.Close()

	shieldF, err := luaEffectFS.Open("effects/shield.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer shieldF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	poisonDef, err := ctx.NewLuaEffect(poisonF)
	if err != nil {
		t.Fatal(err)
	}

	regenDef, err := ctx.NewLuaEffect(regenF)
	if err != nil {
		t.Fatal(err)
	}

	shieldDef, err := ctx.NewLuaEffect(shieldF)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	target.SetEffect(shieldDef, target, 1)
	target.SetEffect(poisonDef, target, 1)
	target.SetEffect(regenDef, target, 1)
	hp, _ := target.GetStat(hpStatDef)
	hp.SetBase(10)

	ctx.Tick()

	hp, _ = target.GetStat(hpStatDef)
	if hp.GetValue() != 78 {
		t.Fatalf("expected hp to be 78 (shield doubles all base changes: 10->20, poison -2->18 doubled->36, regen +3->39 doubled->78), got %d", hp.GetValue())
	}
}

func TestEffectWithNoEventResponse(t *testing.T) {
	f, err := luaEffectFS.Open("effects/copy_stat.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)
	strStatDef := ctx.NewStatDef("str", 0, 99)

	effectDef, err := ctx.NewLuaEffect(f)
	if err != nil {
		t.Fatal(err)
	}

	source := ctx.NewCombatant()
	srcHp, _ := source.GetStat(hpStatDef)
	srcHp.SetBase(30)
	srcStr, _ := source.GetStat(strStatDef)
	srcStr.SetBase(10)

	target := ctx.NewCombatant()
	tgtHp, _ := target.GetStat(hpStatDef)
	tgtHp.SetBase(20)
	tgtStr, _ := target.GetStat(strStatDef)
	tgtStr.SetBase(3)

	target.SetEffect(effectDef, source, 1)

	ctx.Tick()

	tgtStr, _ = target.GetStat(strStatDef)
	if tgtStr.GetValue() != 10 {
		t.Fatalf("expected str to still be 10 (no tick response), got %d", tgtStr.GetValue())
	}

	srcStr, _ = source.GetStat(strStatDef)
	srcStr.SetBase(20)

	tgtStr, _ = target.GetStat(strStatDef)
	if tgtStr.GetValue() != 10 {
		t.Fatalf("expected str to still be 10 (copy_stat only copies on add), got %d", tgtStr.GetValue())
	}
}

func TestDefendEffectWithReducedDamage(t *testing.T) {
	defendF, err := luaEffectFS.Open("effects/defend.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer defendF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 200)

	effectDef, err := ctx.NewLuaEffect(defendF)
	if err != nil {
		t.Fatal(err)
	}

	target := ctx.NewCombatant()
	tgtHp, _ := target.GetStat(hpStatDef)
	tgtHp.SetBase(30)

	attacker := ctx.NewCombatant()
	attHp, _ := attacker.GetStat(hpStatDef)
	attHp.SetBase(30)

	target.SetEffect(effectDef, attacker, 1)

	tgtHp, _ = target.GetStat(hpStatDef)
	tgtHp.SetBase(20)

	tgtHp, _ = target.GetStat(hpStatDef)
	if tgtHp.GetValue() != 25 {
		t.Fatalf("expected target hp to be 25 (took only half the damage: 30->25 instead of 30->20), got %d", tgtHp.GetValue())
	}
}

func TestMultipleAttacks(t *testing.T) {

	ctx := nbattle.New()

	func() {
		attackF, err := luaEffectFS.Open("effects/attack.lua")
		if err != nil {
			t.Fatal(err)
		}
		defer attackF.Close()
		ctx.NewLuaEffect(attackF)
		ctx.NewStatDef("hp", 0, 200)
		ctx.NewStatDef("str", 0, 200)
		ctx.NewStatDef("def", 0, 200)
	}()

	target := ctx.NewCombatant()
	tgtHp, _ := target.GetStat("hp")
	tgtHp.SetBase(30)
	tgtDef, _ := target.GetStat("def")
	tgtDef.SetBase(1)

	source := ctx.NewCombatant()
	srcHp, _ := source.GetStat("hp")
	srcHp.SetBase(30)
	srcStr, _ := source.GetStat("str")
	srcStr.SetBase(2)

	evtHook := func(evt event.Event) error {
		switch evt := evt.(type) {
		case *event.CombatantEffect:
			effectDef, _ := ctx.GetEffectDef(evt.EffectDefID)
			log.Println(effectDef, evt.Potency)
		}
		return nil
	}
	ctx.HookEvents(evtHook)

	for tgtHp.GetValue() > 0 {
		if err := target.SetEffect("attack", source, 1); err != nil {
			t.Fatal(err)
			return
		}
		ctx.Tick()
	}

	if tgtHp.GetValue() > 0 {
		t.Fatalf("expected target hp to be 22")
	}

}

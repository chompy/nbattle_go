package tests

import (
	"embed"
	"testing"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/internal/lua"
)

//go:embed effects/*.lua
var luaEffectFS embed.FS

func TestLuaEffect(t *testing.T) {

	f, err := luaEffectFS.Open("effects/example.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)

	effectDef, err := lua.NewLuaEffect(ctx, f)
	if err != nil {
		t.Fatal(err)
	}

	if effectDef.GetName() != "example_effect" {
		t.Error("expected effect name to be example_effect")
	}

	combatant := ctx.NewCombatant(1)
	enemy := ctx.NewCombatant(2)

	if err := combatant.AddEffect(effectDef, enemy); err != nil {
		t.Fatal(err)
	}

	if combatant.GetStat(hpStatDef).GetValue() != 25 {
		t.Error("expected effect OnAdd to set combatant hp to 25")
	}

}

func TestAttackDefend(t *testing.T) {

	attackF, err := luaEffectFS.Open("effects/attack.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer attackF.Close()
	defendF, err := luaEffectFS.Open("effects/defend.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer defendF.Close()

	ctx := nbattle.New()
	hpStatDef := ctx.NewStatDef("hp", 0, 99)
	strStatDef := ctx.NewStatDef("str", 0, 99)
	defStatDef := ctx.NewStatDef("def", 0, 99)

	attackEffectDef, err := lua.NewLuaEffect(ctx, attackF)
	if err != nil {
		t.Fatal(err)
	}
	defendEffectDef, err := lua.NewLuaEffect(ctx, defendF)
	if err != nil {
		t.Fatal(err)
	}

	combatant := ctx.NewCombatant(1)
	combatant.GetStat(hpStatDef).SetBase(30)
	combatant.GetStat(strStatDef).SetBase(5)

	enemy := ctx.NewCombatant(2)
	enemy.GetStat(hpStatDef).SetBase(30)
	enemy.GetStat(defStatDef).SetBase(1)

	enemy.AddEffect(defendEffectDef, enemy)
	enemy.AddEffect(attackEffectDef, combatant)

	if enemy.GetStat(hpStatDef).GetValue() != 28 {
		t.Error("expect enemy combatant hp to be 28")
	}

}

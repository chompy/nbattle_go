package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go"
)

func TestStatDef(t *testing.T) {
	ctx := nbattle.New()
	strDef := ctx.NewStatDef("str", 0, 99)
	if strDef.GetType() != nbattle.ObjectTypeStatDef {
		t.Errorf("stat def type should be %d", nbattle.ObjectTypeStatDef)
	}
	if strDef.GetName() != "str" {
		t.Error("stat def name should be str")
	}
	if strDef.GetID() != 1 {
		t.Error("stat def id should be 1")
	}
	if strDef.GetMin() != 0 {
		t.Error("stat def min should be 0")
	}
	if strDef.GetMax() != 99 {
		t.Error("stat def max should be 99")
	}

}

func TestStatBase(t *testing.T) {
	ctx := nbattle.New()
	strDef := ctx.NewStatDef("str", 0, 99)
	combatant := ctx.NewCombatant(1)
	strStat := combatant.GetStat(strDef)

	strStat.SetBase(30)
	if strStat.GetBase() != 30 {
		t.Error("stat base should be 30")
	}
	if strStat.GetValue() != 30 {
		t.Error("stat value should be 30")
	}

	strStat.AddBase(5)
	if strStat.GetBase() != 35 {
		t.Error("stat base should be 35")
	}
	if strStat.GetValue() != 35 {
		t.Error("stat value should be 35")
	}

	strStat.AddBase(-50)
	if strStat.GetBase() != 0 {
		t.Error("stat base should be 0")
	}
	if strStat.GetValue() != 0 {
		t.Error("stat value should be 0")
	}

	strStat.SetBase(101)
	if strStat.GetBase() != 99 {
		t.Error("stat base should be 99")
	}
	if strStat.GetValue() != 99 {
		t.Error("stat value should be 99")
	}

	strStat.AddBase(-1)
	if strStat.GetBase() != 98 {
		t.Error("stat base should be 98")
	}
	if strStat.GetValue() != 98 {
		t.Error("stat value should be 98")
	}
}

func TestStatMod(t *testing.T) {
	ctx := nbattle.New()
	strDef := ctx.NewStatDef("str", 0, 99)

	combatant := ctx.NewCombatant(1)
	enemy := ctx.NewCombatant(2)
	ally := ctx.NewCombatant(2)
	strStat := combatant.GetStat(strDef)

	strStat.SetBase(30)

	strStat.SetMod(enemy, -5)
	if strStat.GetValue() != 25 {
		t.Error("stat value should be 25")
	}
	if strStat.GetBase() != 30 {
		t.Error("stat base should be 30")
	}

	strStat.SetMod(enemy, 0)
	if strStat.GetValue() != 30 {
		t.Error("stat value should be 30")
	}
	if strStat.GetBase() != 30 {
		t.Error("stat base should be 30")
	}

	strStat.SetMod(enemy, -10)
	strStat.SetMod(ally, 25)
	if strStat.GetValue() != 45 {
		t.Error("stat value should be 45")
	}
	if strStat.GetBase() != 30 {
		t.Error("stat base should be 30")
	}

	strStat.SetMod(ally, 30)
	if strStat.GetValue() != 50 {
		t.Error("stat value should be 50")
	}
	if strStat.GetBase() != 30 {
		t.Error("stat base should be 30")
	}

	strStat.SetMod(enemy, 0)
	if strStat.GetValue() != 60 {
		t.Error("stat value should be 60")
	}
	if strStat.GetBase() != 30 {
		t.Error("stat base should be 30")
	}
}

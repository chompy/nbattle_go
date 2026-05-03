package nbattle

import (
	"github.com/chompy/nbattle_go/internal/event"
)

type Context struct {
	idCounter   int
	tick        int
	objects     []Object
	eventHooks  []event.Hook
	effectStack []Effect
}

func New() *Context {
	return &Context{0, 0, make([]Object, 0), make([]event.Hook, 0), make([]Effect, 0)}
}

func (c *Context) newObject() BaseObject {
	c.idCounter++
	return BaseObject{c.idCounter, c}
}

func (c *Context) GetObjectByID(ID int) Object {
	for _, obj := range c.objects {
		if obj.GetID() == ID {
			return obj
		}
	}
	return nil
}

func (c *Context) GetObject(obj any) Object {
	switch obj := obj.(type) {
	case Object:
		return obj
	case int:
		return c.GetObjectByID(obj)
	case string:
		statDef, _ := c.GetStatDefByName(obj)
		if statDef != nil {
			return statDef
		}
		effectDef, _ := c.GetEffectDefByName(obj)
		if effectDef != nil {
			return effectDef
		}
	}
	return nil
}

func (c *Context) Tick() int {
	c.tick++
	c.EmitEvent(&event.Tick{Tick: c.tick})
	return c.tick
}

func (c *Context) GetTick() int {
	return c.tick
}

func (c *Context) NewStatDef(name string, min, max int) *StatDef {
	stafDef := &StatDef{c.newObject(), name, min, max}
	c.objects = append(c.objects, stafDef)
	return stafDef
}

func (c *Context) GetStatDefs() []*StatDef {
	out := make([]*StatDef, 0)
	for _, object := range c.objects {
		statDef, ok := object.(*StatDef)
		if ok {
			out = append(out, statDef)
		}
	}
	return out
}

func (c *Context) GetStatDefByID(ID int) (*StatDef, error) {
	statDefObj := c.GetObjectByID(ID)
	if statDefObj == nil {
		return nil, ErrObjectNotFound
	}
	statDef, ok := statDefObj.(*StatDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return statDef, nil
}

func (c *Context) GetStatDefByName(name string) (*StatDef, error) {
	for _, def := range c.objects {
		statDef, ok := def.(*StatDef)
		if !ok {
			continue
		}
		if statDef.GetName() == name {
			return statDef, nil
		}
	}
	return nil, ErrObjectNotFound
}

func (c *Context) GetStatByID(ID int) (*Stat, error) {
	statObj := c.GetObjectByID(ID)
	if statObj == nil {
		return nil, ErrObjectNotFound
	}
	stat, ok := statObj.(*Stat)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return stat, nil
}

func (c *Context) NewEffectDef(name string, create func() Effect) *EffectDef {
	effectDef := &EffectDef{c.newObject(), name, create}
	c.objects = append(c.objects, effectDef)
	return effectDef
}

func (c *Context) GetEffectDefByID(ID int) (*EffectDef, error) {
	effectDefObj := c.GetObjectByID(ID)
	if effectDefObj == nil {
		return nil, ErrObjectNotFound
	}
	effectDef, ok := effectDefObj.(*EffectDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return effectDef, nil
}

func (c *Context) GetEffectDefByName(name string) (*EffectDef, error) {
	for _, def := range c.objects {
		effectDef, ok := def.(*EffectDef)
		if !ok {
			continue
		}
		if effectDef.GetName() == name {
			return effectDef, nil
		}
	}
	return nil, ErrObjectNotFound
}

func (c *Context) NewCombatant(team int) *Combatant {
	combatant := &Combatant{c.newObject(), team, make([]*Stat, 0), make([]*CombatantEffect, 0)}
	c.objects = append(c.objects, combatant)
	c.EmitEvent(&event.NewCombatant{ID: combatant.GetID(), Team: team})
	return combatant
}

func (c *Context) NewCombatantWithID(ID int, team int) *Combatant {
	combatant, err := c.GetCombatantByID(ID)
	if err == ErrObjectNotFound {
		combatant = c.NewCombatant(team)
		combatant.id = ID
	}
	return combatant
}

func (c *Context) GetCombatants() []*Combatant {
	out := make([]*Combatant, 0)
	for _, object := range c.objects {
		combatant, ok := object.(*Combatant)
		if ok {
			out = append(out, combatant)
		}
	}
	return out
}

func (c *Context) GetCombatantByID(ID int) (*Combatant, error) {
	combatantObj := c.GetObjectByID(ID)
	if combatantObj == nil {
		return nil, ErrObjectNotFound
	}
	combatant, ok := combatantObj.(*Combatant)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return combatant, nil
}

func (c *Context) GetCombatantWithStat(stat *Stat) (*Combatant, error) {
	combatants := c.GetCombatants()
	for _, combatant := range combatants {
		for _, combatantStat := range combatant.GetStats() {
			if combatantStat.GetID() == stat.GetID() {
				return combatant, nil
			}
		}
	}
	return nil, ErrObjectNotFound
}

func (c *Context) EmitEvent(event event.Event) {
	for _, hook := range c.eventHooks {
		hook(event)
	}
	for _, combatant := range c.GetCombatants() {
		combatant.HandleEffectEvent(event)
	}

}

func (c *Context) HookEvents(hook event.Hook) {
	c.eventHooks = append(c.eventHooks, hook)
}

func (c *Context) HandleEvent(e event.Event) error {
	switch ev := e.(type) {
	case *event.Tick:
		for c.GetTick() < ev.Tick {
			c.Tick()
		}
	case *event.NewCombatant:
		c.NewCombatantWithID(ev.ID, ev.Team)

	case *event.AddCombatantStat:
		combatant, err := c.GetCombatantByID(ev.CombatantID)
		if err != nil {
			return err
		}
		stat := combatant.GetStat(ev.StatDefID)
		stat.id = ev.StatID

	case *event.StatBase:
		statObj := c.GetObjectByID(ev.StatID)
		if statObj == nil {
			return ErrObjectNotFound
		}
		stat, ok := statObj.(*Stat)
		if !ok {
			return ErrUnexpectedObjectType
		}
		stat.SetBase(ev.Value)
		return nil

	case *event.StatMod:
		statObj := c.GetObjectByID(ev.StatID)
		if statObj == nil {
			return ErrObjectNotFound
		}
		stat, ok := statObj.(*Stat)
		if !ok {
			return ErrUnexpectedObjectType
		}
		stat.SetMod(ev.SourceID, ev.ModValue)

	case *event.AddCombatantEffect:
		target, err := c.GetCombatantByID(ev.TargetID)
		if err != nil {
			return err
		}
		// source optional
		sourceObj := c.GetObjectByID(ev.SourceID)
		if sourceObj != nil {
			// allow nil if source not provided
			target.AddEffect(ev.EffectDefID, sourceObj)
		} else {
			target.AddEffect(ev.EffectDefID, nil)
		}

	case *event.RemoveCombatantEffect:
		target, err := c.GetCombatantByID(ev.TargetID)
		if err != nil {
			return err
		}
		target.RemoveEffect(ev.EffectDefID)
	}
	return nil
}

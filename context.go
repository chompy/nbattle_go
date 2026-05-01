package nbattle

import "github.com/chompy/nbattle_go/internal/event"

type Context struct {
	idCounter  int
	tick       int
	objects    []Object
	eventHooks []event.Hook
}

func New() *Context {
	return &Context{0, 0, make([]Object, 0), make([]event.Hook, 0)}
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
		return statDef
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

func (c *Context) NewEffectDef(name string, create func() Effect) *EffectDef {
	effectDef := &EffectDef{c.newObject(), name, create}
	c.objects = append(c.objects, effectDef)
	return effectDef
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

func (c *Context) Combatants() []*Combatant {
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

func (c *Context) EmitEvent(event event.Event) {
	for _, hook := range c.eventHooks {
		hook(event)
	}
	for _, combatant := range c.Combatants() {
		combatant.HandleEffectEvent(event)
	}

}

func (c *Context) HookEvents(hook event.Hook) {
	c.eventHooks = append(c.eventHooks, hook)
}

func (c *Context) HandleEvent(e event.Event) error {
	switch e := e.(type) {
	case *event.Tick:
		for c.GetTick() < e.Tick {
			c.Tick()
		}
	case *event.NewCombatant:
		c.NewCombatantWithID(e.ID, e.Team)

		/*
			case *event:
				combatant, statDef, err := c.getCombatantAndStatDef(event)
				if err != nil {
					return err
				}
				stat := combatant.Stat(statDef)
				stat.id = event.GetInt(1)
				return nil

			case EventTypeStatBase:
				stat, err := c.getStatByID(event.GetInt(0))
				if err != nil {
					return err
				}
				stat.SetBase(event.GetInt(1))
			case EventTypeStatMod:
				stat, err := c.getStatByID(event.GetInt(0))
				if err != nil {
					return err
				}
				source := c.getObjectByID(event.GetInt(1))
				if source == nil {
					return ErrObjectNotFound
				}
				stat.Mod(source, event.GetInt(2))

			case EventTypeCombatantEffectAdd:
				target, effectDef, source, err := c.getEffectAddParams(event)
				if err != nil {
					return err
				}
				target.AddEffect(effectDef, source)
			case EventTypeCombatantEffectRemove:
				target, effectDef, err := c.getEffectRemoveParams(event)
				if err != nil {
					return err
				}
				target.RemoveEffect(effectDef)
		*/
	}

	return nil
}

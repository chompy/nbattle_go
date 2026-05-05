package nbattle

import (
	"github.com/chompy/nbattle_go/event"
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
	case float32:
		return c.GetObjectByID(int(obj))
	case float64:
		return c.GetObjectByID(int(obj))
	case string:
		statDef, _ := c.GetStatDefByName(obj)
		if statDef != nil {
			return statDef
		}
		effectDef, _ := c.GetEffectDefByName(obj)
		if effectDef != nil {
			return effectDef
		}
	case map[string]any:
		objID, ok := obj["id"].(int)
		if !ok {
			return nil
		}
		return c.GetObjectByID(objID)
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

func (c *Context) NewCombatant() *Combatant {
	combatant := &Combatant{c.newObject(), false, make([]*Stat, 0), make([]*CombatantEffect, 0)}
	c.objects = append(c.objects, combatant)
	combatant.SetActive(true)
	return combatant
}

func (c *Context) newCombatantWithID(ID int) *Combatant {
	combatant, err := c.GetCombatantByID(ID)
	if err == ErrObjectNotFound {
		combatant := &Combatant{BaseObject{ID, c}, false, make([]*Stat, 0), make([]*CombatantEffect, 0)}
		c.objects = append(c.objects, combatant)
		return combatant
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
			if combatantStat == stat {
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

func (c *Context) HandleEvent(evt event.Event) error {
	switch evt := evt.(type) {
	case *event.Tick:
		for c.GetTick() < evt.Tick {
			c.Tick()
		}
	case *event.CombatantUpdate:
		combatant, err := c.GetCombatantByID(evt.CombatantID)
		if err != nil {
			if err == ErrObjectNotFound {
				combatant = c.newCombatantWithID(evt.CombatantID)
				combatant.SetActive(evt.Active)
				return nil
			}
			return err
		}
		combatant.SetActive(evt.Active)
		return nil

	case *event.CombatantStatBase:
		combatant, err := c.GetCombatantByID(evt.CombatantID)
		if err != nil {
			return err
		}
		stat := combatant.GetStat(evt.StatDefID)
		stat.SetBase(evt.Value)

	case *event.CombatantStatMod:
		combatant, err := c.GetCombatantByID(evt.CombatantID)
		if err != nil {
			return err
		}
		stat := combatant.GetStat(evt.StatDefID)
		stat.SetMod(evt.SourceID, evt.ModValue)
		return nil

	case *event.CombatantEffect:
		target, err := c.GetCombatantByID(evt.TargetID)
		if err != nil {
			return err
		}
		target.SetEffect(evt.EffectDefID, evt.Potency, evt.SourceID)

	}
	return nil
}

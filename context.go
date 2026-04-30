package nbattle

type Context struct {
	idCounter  int
	tick       int
	objects    []object
	eventHooks []EventHook
}

func New() *Context {
	return &Context{0, 0, make([]object, 0), make([]EventHook, 0)}
}

func (c *Context) newObject() objectBase {
	c.idCounter++
	return objectBase{c.idCounter, c}
}

func (c *Context) getObjectByID(ID int) object {
	for _, obj := range c.objects {
		if obj.ID() == ID {
			return obj
		}
	}
	return nil
}

func (c *Context) Tick() int {
	c.tick++
	c.EmitEvent(EventTypeTick, c.tick)
	return c.tick
}

func (c *Context) GetTick() int {
	return c.tick
}

func (c *Context) NewStatDef(min, max int) *StatDef {
	stafDef := &StatDef{c.newObject(), min, max}
	c.objects = append(c.objects, stafDef)
	return stafDef
}

func (c *Context) NewEffectDef(create func() Effect) *EffectDef {
	effectDef := &EffectDef{c.newObject(), create}
	c.objects = append(c.objects, effectDef)
	return effectDef
}

func (c *Context) NewCombatant() *Combatant {
	combatant := &Combatant{c.newObject(), make([]*Stat, 0), make([]*combatantEffect, 0)}
	c.objects = append(c.objects, combatant)
	c.EmitEvent(EventTypeCombatantNew, combatant.ID())
	return combatant
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
	combatantObj := c.getObjectByID(ID)
	if combatantObj == nil {
		return nil, ErrObjectNotFound
	}
	combatant, ok := combatantObj.(*Combatant)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return combatant, nil
}

func (c *Context) EmitEvent(eventType EventType, values ...any) {
	for i := range values {
		switch value := values[i].(type) {
		case object:
			values[i] = value.ID()
		}
	}
	c.idCounter++
	event := &Event{c.idCounter, eventType, c.GetTick(), values}
	for _, hook := range c.eventHooks {
		hook(event)
	}
	for _, combatant := range c.Combatants() {
		combatant.HandleEffectEvent(event)
	}

}

func (c *Context) HookEvents(hook EventHook) {
	c.eventHooks = append(c.eventHooks, hook)
}

func (c *Context) HandleEvent(event *Event) error {
	switch event.Type() {
	case EventTypeTick:
		c.Tick()

	case EventTypeCombatantNew:
		combatantID := event.GetInt(0)
		combatant := c.NewCombatant()
		combatant.id = combatantID

	case EventTypeCombatantStatAdd:
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
	}

	return nil
}

func (c *Context) getCombatantAndStatDef(event *Event) (*Combatant, *StatDef, error) {
	combatantObj := c.getObjectByID(event.GetInt(0))
	statDefObj := c.getObjectByID(event.GetInt(2))
	if combatantObj == nil || statDefObj == nil {
		return nil, nil, ErrObjectNotFound
	}
	combatant, ok1 := combatantObj.(*Combatant)
	statDef, ok2 := statDefObj.(*StatDef)
	if !ok1 || !ok2 {
		return nil, nil, ErrUnexpectedObjectType
	}
	return combatant, statDef, nil
}

func (c *Context) getStatByID(id int) (*Stat, error) {
	obj := c.getObjectByID(id)
	if obj == nil {
		return nil, ErrObjectNotFound
	}
	stat, ok := obj.(*Stat)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return stat, nil
}

func (c *Context) getEffectAddParams(event *Event) (*Combatant, *EffectDef, *Combatant, error) {
	targetObj := c.getObjectByID(event.GetInt(0))
	effectDefObj := c.getObjectByID(event.GetInt(1))
	sourceObj := c.getObjectByID(event.GetInt(2))
	if targetObj == nil || effectDefObj == nil || sourceObj == nil {
		return nil, nil, nil, ErrObjectNotFound
	}
	target, ok1 := targetObj.(*Combatant)
	effectDef, ok2 := effectDefObj.(*EffectDef)
	source, ok3 := sourceObj.(*Combatant)
	if !ok1 || !ok2 || !ok3 {
		return nil, nil, nil, ErrUnexpectedObjectType
	}
	return target, effectDef, source, nil
}

func (c *Context) getEffectRemoveParams(event *Event) (*Combatant, *EffectDef, error) {
	targetObj := c.getObjectByID(event.GetInt(0))
	effectDefObj := c.getObjectByID(event.GetInt(1))
	if targetObj == nil || effectDefObj == nil {
		return nil, nil, ErrObjectNotFound
	}
	target, ok1 := targetObj.(*Combatant)
	effectDef, ok2 := effectDefObj.(*EffectDef)
	if !ok1 || !ok2 {
		return nil, nil, ErrUnexpectedObjectType
	}
	return target, effectDef, nil
}

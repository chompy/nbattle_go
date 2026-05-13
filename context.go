package nbattle

import (
	"log/slog"
	"os"
	"slices"

	"github.com/chompy/nbattle_go/event"
)

// Context is the main object of NBattle.
type Context struct {
	idCounter   int
	tick        int
	objects     []Object
	eventHooks  []event.Hook
	effectStack []Effect
	flags       map[string]uint64
	flagCounter uint64
	log         *slog.Logger
}

// New creates a new NBattle context.
func New() *Context {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("New NBattle context created.")
	return &Context{
		idCounter:   0,
		tick:        0,
		objects:     make([]Object, 0),
		eventHooks:  make([]event.Hook, 0),
		effectStack: make([]Effect, 0),
		flags:       make(map[string]uint64),
		flagCounter: 1,
		log:         logger,
	}
}

func (c *Context) newObject() BaseObject {
	c.idCounter++
	return BaseObject{c.idCounter, c}
}

// GetObjectByID retrieves an object by its ID.
func (c *Context) GetObjectByID(ID int) (Object, error) {
	for _, obj := range c.objects {
		if obj.GetID() == ID {
			return obj, nil
		}
	}
	return nil, ErrObjectNotFound
}

// GetObject retrieves an object from an unknown type value (ID, name, or map containing ID).
func (c *Context) GetObject(obj any) (Object, error) {
	switch obj := obj.(type) {
	case Object:
		return obj, nil
	case int:
		return c.GetObjectByID(obj)
	case float32:
		return c.GetObjectByID(int(obj))
	case float64:
		return c.GetObjectByID(int(obj))
	case string:
		return c.GetObjectByName(obj)
	case map[string]any:
		objID, ok := obj["id"].(int)
		if !ok {
			return nil, ErrUnexpectedObjectType
		}
		return c.GetObjectByID(objID)
	}
	return nil, ErrUnexpectedObjectType
}

// GetObjectByName retrieves an object by its name.
func (c *Context) GetObjectByName(name string) (Object, error) {
	for _, object := range c.objects {
		switch object := object.(type) {
		case *StatDef:
			if object.GetName() == name {
				return object, nil
			}
		case *EffectDef:
			if object.GetName() == name {
				return object, nil
			}
		case *TriggerDef:
			if object.GetName() == name {
				return object, nil
			}
		}
	}
	return nil, ErrObjectNotFound
}

// NewStatDef creates a new stat definition.
func (c *Context) NewStatDef(name string, min, max int) *StatDef {
	stafDef := &StatDef{c.newObject(), name, min, max}
	c.objects = append(c.objects, stafDef)
	c.log.Debug("New stat def created.", "object", stafDef)
	return stafDef
}

// GetStatDefs retrieves all stat definitions.
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

// GetStatDefByID retrieves a stat definition by its ID.
func (c *Context) GetStatDefByID(ID int) (*StatDef, error) {
	statDefObj, err := c.GetObjectByID(ID)
	if err != nil {
		return nil, err
	}
	statDef, ok := statDefObj.(*StatDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return statDef, nil
}

// GetStatDefByName retrieves a stat definition by its name.
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

// NewEffectDef creates a new effect definition.
func (c *Context) NewEffectDef(name string, create func() Effect) *EffectDef {
	effectDef := &EffectDef{c.newObject(), name, create}
	c.objects = append(c.objects, effectDef)
	c.log.Debug("New effect def created.", "object", effectDef)
	return effectDef
}

// GetEffectDefByID retrieves an effect definition by its ID.
func (c *Context) GetEffectDefByID(ID int) (*EffectDef, error) {
	effectDefObj, err := c.GetObjectByID(ID)
	if err != nil {
		return nil, err
	}
	effectDef, ok := effectDefObj.(*EffectDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return effectDef, nil
}

// GetEffectDefByName retrieves an effect definition by its name.
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

// NewTriggerDef creates a new trigger definition.
func (c *Context) NewTriggerDef(name string) *TriggerDef {
	triggerDef := &TriggerDef{c.newObject(), name}
	c.objects = append(c.objects, triggerDef)
	c.log.Debug("New trigger def created.", "object", triggerDef)
	return triggerDef
}

// GetTriggerDefByID retrieves a trigger definition by its ID.
func (c *Context) GetTriggerDefByID(ID int) (*TriggerDef, error) {
	triggerDefObj, err := c.GetObjectByID(ID)
	if err != nil {
		return nil, err
	}
	triggerDef, ok := triggerDefObj.(*TriggerDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return triggerDef, nil
}

// GetTriggerDefByName retrieves a trigger definition by its name.
func (c *Context) GetTriggerDefByName(name string) (*TriggerDef, error) {
	for _, def := range c.objects {
		triggerDef, ok := def.(*TriggerDef)
		if !ok {
			continue
		}
		if triggerDef.GetName() == name {
			return triggerDef, nil
		}
	}
	return nil, ErrObjectNotFound
}

// NewCombatant creates a new combatant.
func (c *Context) NewCombatant() *Combatant {
	combatant := &Combatant{c.newObject(), false, make([]*Stat, 0), make([]*CombatantEffect, 0), 0}
	c.objects = append(c.objects, combatant)
	combatant.SetActive(true)
	c.log.Debug("New combatant created.", "object", combatant)
	return combatant
}

func (c *Context) newCombatantWithID(ID int) *Combatant {
	combatant, err := c.GetCombatantByID(ID)
	if err == ErrObjectNotFound {
		combatant := &Combatant{BaseObject{ID, c}, false, make([]*Stat, 0), make([]*CombatantEffect, 0), 0}
		c.objects = append(c.objects, combatant)
		c.log.Debug("New combatant created.", "object", combatant)
		return combatant
	}
	return combatant
}

// GetCombatants retrieves all combatants.
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

// GetCombatantByID retrieves a combatant by its ID.
func (c *Context) GetCombatantByID(ID int) (*Combatant, error) {
	combatantObj, err := c.GetObjectByID(ID)
	if err != nil {
		return nil, err
	}
	combatant, ok := combatantObj.(*Combatant)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return combatant, nil
}

// GetCombatantWithStat retrieves a combatant by a stat that has been assigned to it.
func (c *Context) GetCombatantWithStat(stat *Stat) (*Combatant, error) {
	combatants := c.GetCombatants()
	for _, combatant := range combatants {
		if slices.Contains(combatant.GetStats(), stat) {
			return combatant, nil
		}
	}
	return nil, ErrObjectNotFound
}

// NewFlag creates a new flag.
func (c *Context) NewFlag(name string) uint64 {
	flag := c.flagCounter
	c.flags[name] = flag
	c.flagCounter <<= 1
	return flag
}

// GetFlags retrieves all flags.
func (c *Context) GetFlags() map[string]uint64 {
	return c.flags
}

// GetFlagByName retrieves a flag by its name.
func (c *Context) GetFlagByName(name string) uint64 {
	return c.flags[name]
}

// Tick advances the tick counter and emit tick event.
func (c *Context) Tick() int {
	// clean up effects that have been removed
	for _, combatant := range c.GetCombatants() {
		combatant.processEffectRemovals()
	}

	// emit tick event
	c.tick++
	c.log.Debug("Next tick.", "tick", c.tick)
	c.EmitEvent(&event.Tick{Tick: c.tick})
	return c.tick
}

// GetTick retrieves the current tick.
func (c *Context) GetTick() int {
	return c.tick
}

// EmitEvent sends an event to all hooks and active combatant effects.
func (c *Context) EmitEvent(event event.Event) error {
	c.log.Debug("Emit event.", "event", event.Type())
	for _, hook := range c.eventHooks {
		if err := hook(event); err != nil {
			return err
		}
	}
	for _, combatant := range c.GetCombatants() {
		if err := combatant.HandleEffectEvent(event); err != nil {
			return err
		}
	}
	return nil
}

// HookEvents adds a new event hook.
func (c *Context) HookEvents(hook event.Hook) {
	c.eventHooks = append(c.eventHooks, hook)
}

// HandleEvent processes an event from another context.
func (c *Context) HandleEvent(evt event.Event) error {
	c.log.Debug("Handle event.", "event", evt.Type())
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
		combatant.flags = evt.Flags
		return nil

	case *event.CombatantStatBase:
		combatant, err := c.GetCombatantByID(evt.CombatantID)
		if err != nil {
			return err
		}
		stat, err := combatant.GetStat(evt.StatDefID)
		if err != nil {
			return err
		}
		stat.SetBase(evt.Value)

	case *event.CombatantStatMod:
		combatant, err := c.GetCombatantByID(evt.CombatantID)
		if err != nil {
			return err
		}
		stat, err := combatant.GetStat(evt.StatDefID)
		if err != nil {
			return err
		}
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

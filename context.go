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

// GetObjecGetObjectByIDAndTypetByID retrieves an object by its ID and type.
func (c *Context) GetObjectByIDAndType(ID int, objType ObjectType) (Object, error) {
	for _, obj := range c.objects {
		if obj.GetID() == ID && (objType == ObjectTypeUnknown || objType == obj.GetType()) {
			return obj, nil
		}
	}
	c.log.Error("Unable to find object with ID and type.", "id", ID, "type", objType)
	return nil, ErrObjectNotFound
}

// GetObjectByID retrieves an object by its ID.
func (c *Context) GetObjectByID(ID int) (Object, error) {
	return c.GetObjectByIDAndType(ID, ObjectTypeUnknown)
}

// GetObject retrieves an object from an unknown type value (ID, name, or map containing ID).
func (c *Context) GetObjectByType(obj any, objType ObjectType) (Object, error) {
	switch obj := obj.(type) {
	case Object:
		if objType == ObjectTypeUnknown || obj.GetType() == objType {
			return obj, nil
		}
	case int:
		return c.GetObjectByIDAndType(obj, objType)
	case float32:
		return c.GetObjectByIDAndType(int(obj), objType)
	case float64:
		return c.GetObjectByIDAndType(int(obj), objType)
	case string:
		return c.GetObjectByNameAndType(obj, objType)
	case map[string]any:
		objID, ok := obj["id"].(int)
		if !ok {
			return nil, ErrUnexpectedObjectType
		}
		return c.GetObjectByIDAndType(objID, objType)
	}
	c.log.Error("Unexpected object type.", "object", obj, "type", objType)
	return nil, ErrUnexpectedObjectType
}

// GetObject retrieves an object from an unknown type value (ID, name, or map containing ID).
func (c *Context) GetObject(obj any) (Object, error) {
	return c.GetObjectByType(obj, ObjectTypeUnknown)
}

// GetObjectByNameAndType retrieves an object by its name and type.
func (c *Context) GetObjectByNameAndType(name string, objType ObjectType) (Object, error) {
	for _, object := range c.objects {
		if objType == ObjectTypeUnknown || object.GetType() == objType {
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
	}
	c.log.Error("Object with name and type not found.", "name", name, "type", objType)
	return nil, ErrObjectNotFound
}

// GetObjectByName retrieves an object by its name.
func (c *Context) GetObjectByName(name string) (Object, error) {
	return c.GetObjectByNameAndType(name, ObjectTypeUnknown)
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

// GetStatDef retrieves a stat definition from an unknown variable.
func (c *Context) GetStatDef(obj any) (*StatDef, error) {
	statDefObj, err := c.GetObjectByType(obj, ObjectTypeStatDef)
	if err != nil {
		return nil, err
	}
	statDef, ok := statDefObj.(*StatDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return statDef, nil
}

// NewEffectDef creates a new effect definition.
func (c *Context) NewEffectDef(name string, create func() Effect) *EffectDef {
	effectDef := &EffectDef{c.newObject(), name, create}
	c.objects = append(c.objects, effectDef)
	c.log.Debug("New effect def created.", "object", effectDef)
	return effectDef
}

// GetEffectDef retrieves an effect definition from an unknown variable.
func (c *Context) GetEffectDef(obj any) (*EffectDef, error) {
	effectDefObj, err := c.GetObjectByType(obj, ObjectTypeEffectDef)
	if err != nil {
		return nil, err
	}
	effectDef, ok := effectDefObj.(*EffectDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return effectDef, nil
}

// NewTriggerDef creates a new trigger definition.
func (c *Context) NewTriggerDef(name string) *TriggerDef {
	triggerDef := &TriggerDef{c.newObject(), name}
	c.objects = append(c.objects, triggerDef)
	c.log.Debug("New trigger def created.", "object", triggerDef)
	return triggerDef
}

// GetTriggerDef retrieves a trigger definition from an unknown variable.
func (c *Context) GetTriggerDef(obj any) (*TriggerDef, error) {
	triggerDefObj, err := c.GetObjectByType(obj, ObjectTypeTriggerDef)
	if err != nil {
		return nil, err
	}
	triggerDef, ok := triggerDefObj.(*TriggerDef)
	if !ok {
		return nil, ErrUnexpectedObjectType
	}
	return triggerDef, nil
}

// NewCombatant creates a new combatant.
func (c *Context) NewCombatant() *Combatant {
	combatant := &Combatant{c.newObject(), false, make([]*Stat, 0), make([]*combatantEffect, 0), 0}
	c.objects = append(c.objects, combatant)
	combatant.SetActive(true)
	c.log.Debug("New combatant created.", "object", combatant)
	return combatant
}

func (c *Context) newCombatantWithID(ID int) *Combatant {
	combatant, err := c.GetCombatant(ID)
	if err == ErrObjectNotFound {
		combatant := &Combatant{BaseObject{ID, c}, false, make([]*Stat, 0), make([]*combatantEffect, 0), 0}
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

// GetCombatant retrieves a combatant from an unknown variable.
func (c *Context) GetCombatant(obj any) (*Combatant, error) {
	combatantObj, err := c.GetObjectByType(obj, ObjectTypeCombatant)
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
	for _, object := range c.objects {
		combatant, ok := object.(*Combatant)
		if ok && slices.Contains(combatant.GetStats(), stat) {
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
	for _, combatant := range c.GetCombatants() {
		if err := combatant.processEvent(event); err != nil {
			return err
		}
	}
	for _, hook := range c.eventHooks {
		if err := hook(event); err != nil {
			return err
		}
	}
	return nil
}

// EmitTrigger emits a trigger event from the given source.
func (c *Context) EmitTrigger(triggerDefObj any, sourceObj any) error {
	triggerDef, err := c.GetTriggerDef(triggerDefObj)
	if err != nil {
		return err
	}
	source, err := c.GetObject(sourceObj)
	if err != nil {
		return err
	}
	return c.EmitEvent(&event.Trigger{
		TriggerDefID: triggerDef.GetID(),
		SourceID:     source.GetID(),
	})
}

// HookEvents adds a new event hook.
func (c *Context) HookEvents(hook event.Hook) {
	c.eventHooks = append(c.eventHooks, hook)
}

// ProcessEvent processes an event from another context.
func (c *Context) ProcessEvent(evt event.Event) error {
	c.log.Debug("Process event.", "event", evt.Type())
	switch evt := evt.(type) {
	case *event.Tick:
		for c.GetTick() < evt.Tick {
			c.Tick()
		}
	case *event.CombatantUpdate:
		combatant, err := c.GetCombatant(evt.CombatantID)
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
		combatant, err := c.GetCombatant(evt.CombatantID)
		if err != nil {
			return err
		}
		stat, err := combatant.GetStat(evt.StatDefID)
		if err != nil {
			return err
		}
		stat.SetBase(evt.Value)

	case *event.CombatantStatMod:
		combatant, err := c.GetCombatant(evt.CombatantID)
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
		target, err := c.GetCombatant(evt.TargetID)
		if err != nil {
			return err
		}
		target.SetEffect(evt.EffectDefID, evt.Potency, evt.SourceID)
	}

	return nil
}

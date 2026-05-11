package tests

import (
	"testing"

	nbattle "github.com/chompy/nbattle_go"
	"github.com/chompy/nbattle_go/event"
)

func TestTriggerDefCreation(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("my_trigger")

	if triggerDef.GetType() != nbattle.ObjectTypeTriggerDef {
		t.Errorf("expected type ObjectTypeTriggerDef, got %d", triggerDef.GetType())
	}

	if triggerDef.GetName() != "my_trigger" {
		t.Errorf("expected name 'my_trigger', got '%s'", triggerDef.GetName())
	}

	if triggerDef.GetID() != 1 {
		t.Errorf("expected id 1, got %d", triggerDef.GetID())
	}
}

func TestTriggerDefGetByID(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("my_trigger")

	retrieved, err := ctx.GetTriggerDefByID(triggerDef.GetID())
	if err != nil {
		t.Fatal(err)
	}

	if retrieved != triggerDef {
		t.Fatal("expected to get the same trigger def")
	}

	_, err = ctx.GetTriggerDefByID(999)
	if err != nbattle.ErrObjectNotFound {
		t.Errorf("expected ErrObjectNotFound, got %v", err)
	}
}

func TestTriggerDefGetByName(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("my_trigger")

	retrieved, err := ctx.GetTriggerDefByName("my_trigger")
	if err != nil {
		t.Fatal(err)
	}

	if retrieved != triggerDef {
		t.Fatal("expected to get the same trigger def")
	}

	_, err = ctx.GetTriggerDefByName("nonexistent")
	if err != nbattle.ErrObjectNotFound {
		t.Errorf("expected ErrObjectNotFound, got %v", err)
	}
}

func TestTriggerDefWrongType(t *testing.T) {
	ctx := nbattle.New()
	ctx.NewStatDef("hp", 0, 99)

	_, err := ctx.GetTriggerDefByID(1)
	if err != nbattle.ErrUnexpectedObjectType {
		t.Errorf("expected ErrUnexpectedObjectType, got %v", err)
	}
}

func TestTriggerEventSerialization(t *testing.T) {
	evt := &event.Trigger{
		TriggerDefID:   5,
		EffectDefID:    3,
		EffectTargetID: 10,
		EffectSourceID: 7,
		EffectPotency:  2,
	}

	data, err := evt.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if evt.Type() != event.TriggerEvent {
		t.Errorf("expected type TriggerEvent, got %d", evt.Type())
	}

	deserialized := &event.Trigger{}
	err = deserialized.Deserialize(data)
	if err != nil {
		t.Fatal(err)
	}

	if deserialized.TriggerDefID != 5 {
		t.Errorf("expected TriggerDefID 5, got %d", deserialized.TriggerDefID)
	}
	if deserialized.EffectDefID != 3 {
		t.Errorf("expected EffectDefID 3, got %d", deserialized.EffectDefID)
	}
	if deserialized.EffectTargetID != 10 {
		t.Errorf("expected EffectTargetID 10, got %d", deserialized.EffectTargetID)
	}
	if deserialized.EffectSourceID != 7 {
		t.Errorf("expected EffectSourceID 7, got %d", deserialized.EffectSourceID)
	}
	if deserialized.EffectPotency != 2 {
		t.Errorf("expected EffectPotency 2, got %d", deserialized.EffectPotency)
	}
}

func TestTriggerEventDeserializationWrongType(t *testing.T) {
	tickEvt := &event.Tick{Tick: 1}
	data, err := tickEvt.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	triggerEvt := &event.Trigger{}
	err = triggerEvt.Deserialize(data)
	if err != event.ErrDeserializeWrongType {
		t.Errorf("expected ErrDeserializeWrongType, got %v", err)
	}
}

type TriggerCaptureEffect struct {
	TriggerDef     *nbattle.TriggerDef
	CapturedEvents []*event.Trigger
}

func (e *TriggerCaptureEffect) OnAdd(ctx *nbattle.EffectCtx) {
}

func (e *TriggerCaptureEffect) OnRemove(ctx *nbattle.EffectCtx) {
}

func (e *TriggerCaptureEffect) OnEvent(ctx *nbattle.EffectCtx, evt event.Event) {
	if triggerEvt, ok := evt.(*event.Trigger); ok {
		e.CapturedEvents = append(e.CapturedEvents, triggerEvt)
	}
}

func TestTriggerEmitViaEffectCtx(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("damage_trigger")

	captureEffect := &TriggerCaptureEffect{TriggerDef: triggerDef}
	captureEffectDef := ctx.NewEffectDef("capture", func() nbattle.Effect {
		return captureEffect
	})

	target := ctx.NewCombatant()
	target.SetEffect(captureEffectDef, 1, nil)

	source := ctx.NewCombatant()

	emitEffect := &EmitTriggerEffect{TriggerDef: triggerDef}
	emitEffectDef := ctx.NewEffectDef("emitter", func() nbattle.Effect {
		return emitEffect
	})

	source.SetEffect(emitEffectDef, 1, nil)

	if len(captureEffect.CapturedEvents) != 1 {
		t.Fatalf("expected 1 captured trigger event, got %d", len(captureEffect.CapturedEvents))
	}

	evt := captureEffect.CapturedEvents[0]
	if evt.TriggerDefID != triggerDef.GetID() {
		t.Errorf("expected TriggerDefID %d, got %d", triggerDef.GetID(), evt.TriggerDefID)
	}
	if evt.EffectDefID != emitEffectDef.GetID() {
		t.Errorf("expected EffectDefID %d, got %d", emitEffectDef.GetID(), evt.EffectDefID)
	}
	if evt.EffectTargetID != source.GetID() {
		t.Errorf("expected EffectTargetID %d, got %d", source.GetID(), evt.EffectTargetID)
	}
	if evt.EffectPotency != 1 {
		t.Errorf("expected EffectPotency 1, got %d", evt.EffectPotency)
	}
}

type EmitTriggerEffect struct {
	TriggerDef *nbattle.TriggerDef
}

func (e *EmitTriggerEffect) OnAdd(ctx *nbattle.EffectCtx) {
	ctx.EmitTrigger(e.TriggerDef)
}

func (e *EmitTriggerEffect) OnRemove(ctx *nbattle.EffectCtx) {
}

func (e *EmitTriggerEffect) OnEvent(ctx *nbattle.EffectCtx, evt event.Event) {
}

func TestTriggerEmitViaTriggerDef(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("test_trigger")

	captureEffect := &TriggerCaptureEffect{}
	captureEffectDef := ctx.NewEffectDef("capture", func() nbattle.Effect {
		return captureEffect
	})

	target := ctx.NewCombatant()
	target.SetEffect(captureEffectDef, 1, nil)

	source := ctx.NewCombatant()

	effectCtx := &nbattle.EffectCtx{
		Ctx:     ctx,
		Def:     ctx.NewEffectDef("dummy", func() nbattle.Effect { return &EmitTriggerEffect{} }),
		Potency: 3,
		Target:  source,
		Source:  target,
	}

	triggerDef.EmitEvent(effectCtx)

	if len(captureEffect.CapturedEvents) != 1 {
		t.Fatalf("expected 1 captured trigger event, got %d", len(captureEffect.CapturedEvents))
	}

	evt := captureEffect.CapturedEvents[0]
	if evt.TriggerDefID != triggerDef.GetID() {
		t.Errorf("expected TriggerDefID %d, got %d", triggerDef.GetID(), evt.TriggerDefID)
	}
	if evt.EffectPotency != 3 {
		t.Errorf("expected EffectPotency 3, got %d", evt.EffectPotency)
	}
	if evt.EffectSourceID != target.GetID() {
		t.Errorf("expected EffectSourceID %d, got %d", target.GetID(), evt.EffectSourceID)
	}
}

func TestTriggerWithNilSource(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("test_trigger")

	captureEffect := &TriggerCaptureEffect{}
	captureEffectDef := ctx.NewEffectDef("capture", func() nbattle.Effect {
		return captureEffect
	})

	target := ctx.NewCombatant()
	target.SetEffect(captureEffectDef, 1, nil)

	effectCtx := &nbattle.EffectCtx{
		Ctx:     ctx,
		Def:     ctx.NewEffectDef("dummy", func() nbattle.Effect { return &EmitTriggerEffect{} }),
		Potency: 1,
		Target:  target,
		Source:  nil,
	}

	triggerDef.EmitEvent(effectCtx)

	if len(captureEffect.CapturedEvents) != 1 {
		t.Fatalf("expected 1 captured trigger event, got %d", len(captureEffect.CapturedEvents))
	}

	evt := captureEffect.CapturedEvents[0]
	if evt.EffectSourceID != 0 {
		t.Errorf("expected EffectSourceID 0 for nil source, got %d", evt.EffectSourceID)
	}
}

func TestTriggerMultipleEffects(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("multi_trigger")

	capture1 := &TriggerCaptureEffect{}
	capture2 := &TriggerCaptureEffect{}

	ctx.NewEffectDef("capture1", func() nbattle.Effect { return capture1 })
	ctx.NewEffectDef("capture2", func() nbattle.Effect { return capture2 })

	target := ctx.NewCombatant()
	cap1Def, _ := ctx.GetEffectDefByName("capture1")
	cap2Def, _ := ctx.GetEffectDefByName("capture2")
	target.SetEffect(cap1Def, 1, nil)
	target.SetEffect(cap2Def, 1, nil)

	source := ctx.NewCombatant()
	sourceHp, _ := source.GetStat(ctx.NewStatDef("hp", 0, 99))
	sourceHp.SetBase(50)

	effectCtx := &nbattle.EffectCtx{
		Ctx:     ctx,
		Def:     ctx.NewEffectDef("dummy", func() nbattle.Effect { return &EmitTriggerEffect{} }),
		Potency: 5,
		Target:  source,
		Source:  target,
	}

	triggerDef.EmitEvent(effectCtx)

	if len(capture1.CapturedEvents) != 1 {
		t.Errorf("expected capture1 to have 1 event, got %d", len(capture1.CapturedEvents))
	}
	if len(capture2.CapturedEvents) != 1 {
		t.Errorf("expected capture2 to have 1 event, got %d", len(capture2.CapturedEvents))
	}
}

func TestTriggerEmitTriggerOnNotFound(t *testing.T) {
	ctx := nbattle.New()
	effectCtx := &nbattle.EffectCtx{
		Ctx:     ctx,
		Def:     ctx.NewEffectDef("dummy", func() nbattle.Effect { return &EmitTriggerEffect{} }),
		Potency: 1,
		Target:  ctx.NewCombatant(),
		Source:  nil,
	}

	err := effectCtx.EmitTrigger(999)
	if err != nbattle.ErrObjectNotFound {
		t.Errorf("expected ErrObjectNotFound, got %v", err)
	}
}

func TestTriggerEmitTriggerWrongType(t *testing.T) {
	ctx := nbattle.New()
	ctx.NewStatDef("hp", 0, 99)

	effectCtx := &nbattle.EffectCtx{
		Ctx:     ctx,
		Def:     ctx.NewEffectDef("dummy", func() nbattle.Effect { return &EmitTriggerEffect{} }),
		Potency: 1,
		Target:  ctx.NewCombatant(),
		Source:  nil,
	}

	err := effectCtx.EmitTrigger(1)
	if err != nbattle.ErrUnexpectedObjectType {
		t.Errorf("expected ErrUnexpectedObjectType, got %v", err)
	}
}

func TestTriggerDefGetObjectByName(t *testing.T) {
	ctx := nbattle.New()
	ctx.NewStatDef("hp", 0, 99)
	triggerDef := ctx.NewTriggerDef("my_trigger")

	obj, err := ctx.GetObjectByName("my_trigger")
	if err != nil {
		t.Fatal(err)
	}

	if obj != triggerDef {
		t.Fatal("expected to get the trigger def by name")
	}
}

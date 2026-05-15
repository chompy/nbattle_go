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

	retrieved, err := ctx.GetTriggerDef(triggerDef.GetID())
	if err != nil {
		t.Fatal(err)
	}

	if retrieved != triggerDef {
		t.Fatal("expected to get the same trigger def")
	}

	_, err = ctx.GetTriggerDef(999)
	if err != nbattle.ErrObjectNotFound {
		t.Errorf("expected ErrObjectNotFound, got %v", err)
	}
}

func TestTriggerDefGetByName(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("my_trigger")

	retrieved, err := ctx.GetTriggerDef("my_trigger")
	if err != nil {
		t.Fatal(err)
	}

	if retrieved != triggerDef {
		t.Fatal("expected to get the same trigger def")
	}

	_, err = ctx.GetTriggerDef("nonexistent")
	if err != nbattle.ErrObjectNotFound {
		t.Errorf("expected ErrObjectNotFound, got %v", err)
	}
}

func TestTriggerEventSerialization(t *testing.T) {
	evt := &event.Trigger{
		TriggerDefID: 5,
		SourceID:     7,
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
	if deserialized.SourceID != 7 {
		t.Errorf("expected SourceID 7, got %d", deserialized.SourceID)
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

func (e *TriggerCaptureEffect) OnAdd(ctx *nbattle.Context, effectCtx *nbattle.EffectContext) {
}

func (e *TriggerCaptureEffect) OnRemove(ctx *nbattle.Context, effectCtx *nbattle.EffectContext) {
}

func (e *TriggerCaptureEffect) OnEvent(ctx *nbattle.Context, effectCtx *nbattle.EffectContext, evt event.Event) {
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
	target.SetEffect(captureEffectDef, target, 1)

	source := ctx.NewCombatant()

	emitEffect := &EmitTriggerEffect{TriggerDef: triggerDef}
	emitEffectDef := ctx.NewEffectDef("emitter", func() nbattle.Effect {
		return emitEffect
	})

	source.SetEffect(emitEffectDef, source, 1)

	if len(captureEffect.CapturedEvents) != 1 {
		t.Fatalf("expected 1 captured trigger event, got %d", len(captureEffect.CapturedEvents))
	}

	evt := captureEffect.CapturedEvents[0]
	if evt.TriggerDefID != triggerDef.GetID() {
		t.Errorf("expected TriggerDefID %d, got %d", triggerDef.GetID(), evt.TriggerDefID)
	}
	if evt.SourceID != source.GetID() {
		t.Errorf("expected SourceID %d, got %d", source.GetID(), evt.SourceID)
	}
}

type EmitTriggerEffect struct {
	TriggerDef *nbattle.TriggerDef
}

func (e *EmitTriggerEffect) OnAdd(ctx *nbattle.Context, effectCtx *nbattle.EffectContext) {
	ctx.EmitTrigger(e.TriggerDef, effectCtx.Target)
}

func (e *EmitTriggerEffect) OnRemove(ctx *nbattle.Context, effectCtx *nbattle.EffectContext) {
}

func (e *EmitTriggerEffect) OnEvent(ctx *nbattle.Context, effectCtx *nbattle.EffectContext, evt event.Event) {
}

func TestTriggerEmitViaTriggerDef(t *testing.T) {
	ctx := nbattle.New()
	triggerDef := ctx.NewTriggerDef("test_trigger")

	captureEffect := &TriggerCaptureEffect{}
	captureEffectDef := ctx.NewEffectDef("capture", func() nbattle.Effect {
		return captureEffect
	})

	target := ctx.NewCombatant()
	target.SetEffect(captureEffectDef, target, 1)

	source := ctx.NewCombatant()
	ctx.EmitTrigger(triggerDef, source)

	if len(captureEffect.CapturedEvents) != 1 {
		t.Fatalf("expected 1 captured trigger event, got %d", len(captureEffect.CapturedEvents))
	}

	evt := captureEffect.CapturedEvents[0]
	if evt.TriggerDefID != triggerDef.GetID() {
		t.Errorf("expected TriggerDefID %d, got %d", triggerDef.GetID(), evt.TriggerDefID)
	}
	if evt.SourceID != source.GetID() {
		t.Errorf("expected SourceID %d, got %d", source.GetID(), evt.SourceID)
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
	cap1Def, _ := ctx.GetEffectDef("capture1")
	cap2Def, _ := ctx.GetEffectDef("capture2")
	target.SetEffect(cap1Def, target, 1)
	target.SetEffect(cap2Def, target, 1)

	source := ctx.NewCombatant()
	sourceHp, _ := source.GetStat(ctx.NewStatDef("hp", 0, 99))
	sourceHp.SetBase(50)

	ctx.EmitTrigger(triggerDef, source)

	if len(capture1.CapturedEvents) != 1 {
		t.Errorf("expected capture1 to have 1 event, got %d", len(capture1.CapturedEvents))
	}
	if len(capture2.CapturedEvents) != 1 {
		t.Errorf("expected capture2 to have 1 event, got %d", len(capture2.CapturedEvents))
	}
}

func TestTriggerEmitTriggerOnNotFound(t *testing.T) {
	ctx := nbattle.New()
	err := ctx.EmitTrigger(999, 998)
	if err != nbattle.ErrObjectNotFound {
		t.Errorf("expected ErrObjectNotFound, got %v", err)
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

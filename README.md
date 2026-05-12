# NBattle (Go)
**By Nathan Ogden**

Go library for managing RPG combat. It allows for the creation of combatants with stats and effects that can be applied to them. Effects are how combatants interact with one another, this can range from basic attacks to long lived status effects. Effects can be written in either Go or Lua.

This is a work in progress and is currently largely untested.

## Lua Effects

Effects can be defined via Lua scripting.


### Usage

To add a Lua effect use the `NewEffect` function and provide the NBattle context and a reader containing the Lua script.

```go
import (
	nbattle "github.com/chompy/nbattle_go"
)

func main() {
    ctx := nbattle.New()
    f, err := os.Open("script.lua")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    effectDef, err := ctx.NewLuaEffect(ctx, f)
    // ...
}
```


### Scripts

The Lua script at minimum should contain a `Name` function that returns the name of the effect.

```lua
function Name()
    return "attack"
end
```

*All variable and argument types...*

- GLOBALS
  - `ctx` - `Context`
  - `STAT_*` - `StatDef` - Every stat definition will get a global value that is 'STAT_` combined with the stat definition name in all uppercase letters.
  - `FLAG_*` - `number` - Every flag will get a global value that is 'FLAG_' combined with the flag name in all uppercase letters.

- `Object` - Value representing an object. It can be the actual object, the ID number of the object, or the name string of the object.

- `Context`
  - `getTick() <number>` - Function that returns the current tick.
  - `getCombatants() <[]Combatant>` - Function that returns a list of all combatants.
  - `getObject(Object) <Object>` - Function that retrieve an object.

- `StatDef`
  - `type <number>` - Type number of object.
  - `name <string>` - Name of the stat definition.
  - `min <number>` - Minimum value of the stat.
  - `max <number>` - Maximum value of the stat.

- `Stat`
  - `def <StatDef>` - Stat definition for this stat.
  - `getBase() <number>` - Function that returns the current base value of this stat.
  - `setBase(<number>)` - Function that sets the base value of this stat.
  - `addBase(<number>)` - Function that adds the given value to the base value of this stat.
  - `subBase(<number>)` - Function that subtracts the given value to the base value of this stat.
  - `getValue(<number>)` - Function that returns the current value of this stat with modifications.
  - `setMod(<Object>, <number>)` - Function that applies a modification to the stat. Modifications are mapped to an object. Modification value is added to the base value. A value of zero removes the mod.
  - `getCombatant() <Combatant>` - Function that returns the combatant this stat is applied to.

- `Combatant`
  - `getStat(statDef <Object>)` - Retrieve a combatant stat from the stat definition.
  - `setEffect(effectDef <Object>, potency <number>, sourceObject <Object>)` - Apply effect to combatant.
  - `removeEffect(effectDef <Object>)` - Remove effect from combatant. Same as calling `setEffect` with zero potency.
  - `hasEffect(effectDef <Object>)` - Return true if combatant has given effect applied.
  - `setFlag(flag <number>)` - Apply a flag to the combatant.
  - `hasFlag(flag <number>)` - Return true if combatant has the given flag.

- `EffectCtx`
  - `target <Combatant>` - The combatant that is the target of the effect.
  - `source <Object>` - The object that is the source of the effect. This can be nil.
  - `effect <string>` - The name of the effect.
  - `potency <number>` - The potency of the effect.
  - `emitTrigger(<string>)` - Function that emits a trigger event.
  - `remove()` - Function that removes this effect. Same as calling `combatant.setEffect` with zero potency.

- `TickEvent`
  - `tick <number>` - The current tick.

- `CombatantUpdateEvent`
  - `combatant <Combatant>` - The combatant who received the update.
  - `active <bool>` - Whether or not the combatant is currently active.
  - `flags <number>` - The combatant flag value.

- `CombatantStatBaseEvent`
  - `combatant <Combatant>` - The combatant whose base stat changed.
  - `statDef <StatDef>` - The stat definition of the stat that changed.
  - `value <number>` - The new value of the stat.
  - `setValue(value <number>)` - Function that allows the base stat value to be altered as part of the same event.

- `CombatantStatModEvent`
  - `combatant <Combatant>` - The combatant whose stat has been modified.
  - `statDef <StatDef>` - The stat definition of the stat that has been modified.
  - `value <number>` - The modification amount.
  - `setValue(value <number>)` - Function that allows the stat modification value to be altered as part of the same event.

- `CombatantEffectEvent`
  - `target <Combatant>` - The combatant whom the effect is being applied to.
  - `effect <string>` - The name of the effect.
  - `potency <number>` - The potency of the effect. If this is zero then the effect is being removed.
  - `source <Object>` - The source of the effect, can be nil.

- `TriggerEvent`
  - `trigger <string>` - The name of the trigger.
  - `target <Combatant>` - The target of the effect that spawned this trigger.
  - `effect <string>` - The effect name that spawned this trigger.
  - `potency <number>` - The potency of the effect that spawned this trigger.
  - `source <Object>` - The source of the effect that spawned this trigger.

*Event handling functions...*

- `OnAdd(ctx <EffectCtx>)` - Called when the effect is applied to a combatant.
- `OnRemove(ctx <EffectCtx>)` - Called when the effect is removed from a combatant.
- `OnTick(ctx <EffectCtx>, evt <TickEvent>)` - Called whenever combatant advances a tick.
- `OnCombatantUpdate(ctx <EffectCtx>, evt <CombatantUpdateEvent>)` - Called when a new combatant is added or an existing combatant is updated.
- `OnCombatantStatBase(ctx <EffectCtx>, evt <CombatantStatBaseEvent>)` - Called when a stat applied to a combatant has a base value change.
- `OnCombatantStatMod(ctx <EffectCtx>, evt <CombatantStatModEvent>)` - Called when a stat applied to a combatant has a mod value change.
- `OnCombatantEffect(ctx <EffectCtx>, evt <CombatantEffectEvent>)` - Called when an effect is added or removed from a combatant.
- `OnTrigger(ctx <EffectCtx>, evt <TriggerEvent>)` - Called when an effect emits a trigger event.

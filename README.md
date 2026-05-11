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
	nbattlelua "github.com/chompy/nbattle_go/lua"
)

func main() {
    ctx := nbattle.New()
    f, err := os.Open("script.lua")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    effectDef, err := nbattlelua.NewEffect(ctx, f)
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
  - `id` - ID of the combatant object.
  - `type` - Type number of object.
  - `getStat(statDef <Object>)` - Retrieve a combatant stat from the stat definition.
  - `setEffect(effectDef <Object>, potency <number>, sourceObject <Object>)` -  

*Event handling functions...*

- `OnAdd(ctx)` - Called when the effect is applied to a combatant.
- `OnRemove(ctx)` - Called when the effect is removed from a combatant.
- `OnTick(ctx, evt)` - Called whenever combatant advances a tick.
- `OnCombatantUpdate(ctx, evt)` - Called when a new combatant is added or an existing combatant is updated.
- `OnCombatantStatBase(ctx, evt)` - Called when a stat applied to a combatant has a base value change.
- `OnCombatantStatMod(ctx, evt)` - Called when a stat applied to a combatant has a mod value change.
- `OnCombatantEffect(ctx, evt)` - Called when an effect is added or removed from a combatant.
- `OnTrigger(ctx, evt)` - Called when an effect emits a trigger event.

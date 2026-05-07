# NBattle - Lua

This submodule allows effects to be defined from Lua scripts. 


## Usage

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


## Scripts

The Lua script at minimum should contain a `Name` function that returns the name of the effect.

```lua
function Name()
    return "attack"
end
```

...TODO/WIP...

All variable and argument types...

- GLOBALS
  - `ctx` - `Context`
  - `STAT_*` - `StatDef` - Every stat definition will get a global value that is 'STAT_` combined with the stat definition name in all uppercase letters.
  - `FLAG_*` - `number` - Every flag will get a global value that is 'FLAG_' combined with the flag name in all uppercase letters.

- `Object` - Value representing an object. It can be the actual object, the ID number of the object, or the name string of the object.

- `Context`
  - `getTick()` - Function that returns the current tick.
  - `getCombatants()` - Function that returns a list of all combatants.
  - `getObject(Object)` - Function that retrieve an object.

- `StatDef`
  - `type` - `number` - Type number of object.
  - `name` - `string` - Name of the stat definition.
  - `min` - `number` - Minimum value of the stat.
  - `max` - `number` - Maximum value of the stat.

- `Stat`
  - TODO

- `Combatant`
  - `id` - ID of the combatant object.
  - `type` - Type number of object.
  - `getStat(statDef <Object>)` - Retrieve a combatant stat from the stat definition.
  - `setEffect(effectDef <Object>, potency <number>, sourceObject <Object>)` -  

Event handling functions...

- `OnAdd(ctx)` - Called when the effect is applied to a combatant.
- `OnRemove(ctx)` - Called when the effect is removed from a combatant.
- `OnTick(ctx, evt)` - Called whenever combatant advances a tick.
- `OnCombatantUpdate(ctx, evt)` - Called when a new combatant is added or an existing combatant is updated.
- `OnCombatantStatBase(ctx, evt)` - Called when a stat applied to a combatant has a base value change.
- `OnCombatantStatMod(ctx, evt)` - Called when a stat applied to a combatant has a mod value change.
- `OnCombatantEffect(ctx, evt)` - Called when an effect is added or removed from a combatant.
- `OnTrigger(ctx, evt)` - Called when an effect emits a trigger event.

function Name()
    return "shield"
end

function OnCombatantStatBase(ctx, evt)
    if ctx.target.id == evt.combatant.id then
        evt.setValue(evt.value * 2)
    end
end
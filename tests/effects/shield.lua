function Name()
    return "shield"
end

function OnCombatantStatBase(effectCtx, evt)
    if effectCtx.target.id == evt.combatant.id then
        evt.setValue(evt.value * 2)
    end
end
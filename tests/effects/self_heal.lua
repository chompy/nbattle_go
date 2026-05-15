function Name()
    return "self_heal"
end

function OnCombatantStatBase(effectCtx, evt)
    if evt.statDef.name == "hp" and effectCtx.target.id == evt.combatant.id then
        local currentHp = effectCtx.target.getStat("hp").getValue()
        if evt.value < currentHp then
            evt.setValue(evt.value + 1)
        end
    end
end


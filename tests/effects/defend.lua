function Name()
    return "defend"
end

function OnCombatantStatBase(effectCtx, evt)
    if evt.statDef.name == "hp" and effectCtx.target.id == evt.combatant.id then
        local diff = effectCtx.target.getStat("hp").getValue() - evt.value
        if diff > 0 then
            evt.setValue(effectCtx.target.getStat("hp").getValue() - (diff / 2))
        end
    end
end

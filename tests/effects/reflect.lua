function Name()
    return "reflect"
end

function OnCombatantStatBase(effectCtx, evt)
    if evt.statDef.name == "hp" and effectCtx.target.id == evt.combatant.id then
        local currentHp = effectCtx.target.getStat("hp").getValue()
        local diff = currentHp - evt.value
        if diff > 0 and effectCtx.source ~= nil then
            effectCtx.source.getStat("hp").subBase(diff)
        end
    end
end
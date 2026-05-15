function Name()
    return "counter"
end

function OnCombatantStatBase(effectCtx, evt)
    if effectCtx.target.id == evt.combatant.id then
        local currentHp = effectCtx.target.getStat("hp").getValue()
        if evt.value < currentHp and effectCtx.source ~= nil then
            effectCtx.source.getStat("hp").subBase(5)
        end
    end
end


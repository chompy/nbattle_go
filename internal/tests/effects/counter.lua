function Name()
    return "counter"
end

function OnCombatantStatBase(ctx, evt)
    if ctx.target.id == evt.combatant.id then
        local currentHp = ctx.target.getStat("hp").getValue()
        if evt.value < currentHp and ctx.source ~= nil then
            ctx.source.getStat("hp").subBase(5)
        end
    end
end


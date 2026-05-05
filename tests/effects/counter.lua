function Name()
    return "counter"
end

function OnCombatantStatBase(ctx, evt)
    if ctx.target.id == evt.combatant.id then
        local currentHp = ctx.target.getStat("hp").get()
        if evt.value < currentHp and ctx.source ~= nil then
            ctx.source.getStat("hp").subtract(5)
        end
    end
end


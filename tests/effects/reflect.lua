function Name()
    return "reflect"
end

function OnCombatantStatBase(ctx, evt)
    if evt.statDef.name == "hp" and ctx.target.id == evt.combatant.id then
        local currentHp = ctx.target.getStat("hp").getValue()
        local diff = currentHp - evt.value
        if diff > 0 and ctx.source ~= nil then
            ctx.source.getStat("hp").subBase(diff)
        end
    end
end
function Name()
    return "defend"
end

function OnCombatantStatBase(ctx, evt)
    if evt.statDef.name == "hp" and ctx.target.id == evt.combatant.id then
        local diff = ctx.target.getStat("hp").get() - evt.value
        if diff > 0 then
            evt.setValue(ctx.target.getStat("hp").get() - (diff / 2))
        end
    end
end

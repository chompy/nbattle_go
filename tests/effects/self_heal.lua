function Name()
    return "self_heal"
end

function OnCombatantStatBase(ctx, evt)
    if evt.statDef.name == "hp" and ctx.target.id == evt.combatant.id then
        local currentHp = ctx.target.getStat("hp").get()
        if evt.value < currentHp then
            evt.setValue(evt.value + 1)
        end
    end
end


function Name()
    return "trigger_effect"
end

function OnCombatantStatBase(ctx, evt)
    if evt.statDef.name == "hp" and ctx.target.id == evt.combatant.id then
        if evt.value <= 0 then
            ctx.target.setEffect("buff", 1, ctx.source)
        end
    end
end

function Name()
    return "trigger_effect"
end

function OnCombatantStatBase(effectCtx, evt)
    if evt.statDef.name == "hp" and effectCtx.target.id == evt.combatant.id then
        if evt.value <= 0 then
            effectCtx.target.setEffect("buff", effectCtx.source, 1)
        end
    end
end

function Name()
    return "trigger_effect"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == COMBATANT_STAT_BASE) then
        if (event.stat.getCombatant().id == target.id and event.stat.def.name == "hp") then
            if event.value <= 0 then
                target.setEffect("buff", 1, source)
            end
        end
    end
end

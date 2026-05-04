function Name()
    return "shield"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == COMBATANT_STAT_BASE) then
        if (event.stat.getCombatant().id == target.id) then
            event.setValue(event.value * 2)
        end
    end
end

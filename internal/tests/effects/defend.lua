function Name()
    return "defend"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == COMBATANT_STAT_BASE) then
        if (event.stat.getCombatant().id == target.id) then
            local diff = target.getStat("hp").get() - event.value
            if diff > 0 then
                event.setValue(target.getStat("hp").get() - (diff / 2))
            end
        end
    end
end

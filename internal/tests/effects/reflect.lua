function Name()
    return "reflect"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == COMBATANT_STAT_BASE) then
        if (event.stat.getCombatant().id == target.id and event.stat.def.name == "hp") then
            local currentHp = target.getStat("hp").get()
            local diff = currentHp - event.value
            if diff > 0 and source ~= nil then
                source.getStat("hp").add(-diff)
            end
        end
    end
end

function Name()
    return "counter"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == COMBATANT_STAT_BASE) then
        if (event.stat.getCombatant().id == target.id and event.stat.def.name == "hp") then
            local currentHp = target.getStat("hp").get()
            if event.value < currentHp and source ~= nil then
                source.getStat("hp").add(-5)
            end
        end
    end
end

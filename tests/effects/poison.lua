function Name()
    return "poison"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == TICK) then
        target.getStat("hp").add(-2)
    end
end

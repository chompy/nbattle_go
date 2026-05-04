function Name()
    return "regenerate"
end

function OnAdd(target, source)
end

function OnRemove()
end

function OnEvent(event, target, source)
    if (event.type == TICK) then
        target.getStat("hp").add(3)
    end
end

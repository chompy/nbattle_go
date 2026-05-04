function Name()
    return "buff"
end

function OnAdd(target, source)
    target.getStat("str").add(10)
end

function OnRemove()
end

function OnEvent(event, target, source)
end

function Name()
    return "example_effect"
end

function OnAdd(target, source)
    target.getStat("hp").set(25)
end

function OnRemove()
end

function OnEvent(event)
end
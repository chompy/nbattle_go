function Name()
    return "attack"
end

function OnAdd(target, source)
    target.getStat("hp").add(-1 * (source.getStat("str").get() - target.getStat("def").get()))
end

function OnRemove()
end

function OnEvent(event)
end
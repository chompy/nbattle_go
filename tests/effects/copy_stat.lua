function Name()
    return "copy_stat"
end

function OnAdd(target, source)
    local strStat = target.getStat("str")
    local sourceStr = source.getStat("str").get()
    strStat.set(sourceStr)
end

function OnRemove()
end

function OnEvent(event, target, source)
end

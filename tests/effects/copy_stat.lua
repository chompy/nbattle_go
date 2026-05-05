function Name()
    return "copy_stat"
end

function OnAdd(ctx)
    local strStat = ctx.target.getStat("str")
    local sourceStr = ctx.source.getStat("str").get()
    strStat.set(sourceStr)
    --ctx.target.removeEffect(ctx.effect)
end
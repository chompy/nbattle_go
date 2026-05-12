function Name()
    return "copy_stat"
end

function OnAdd(ctx)
    local strStat = ctx.target.getStat("str")
    local sourceStr = ctx.source.getStat("str").getValue()
    strStat.setBase(sourceStr)
    ctx.target.removeEffect(ctx.effect)
end
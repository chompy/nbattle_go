function Name()
    return "copy_stat"
end

function OnAdd(effectCtx)
    local strStat = effectCtx.target.getStat("str")
    local sourceStr = effectCtx.source.getStat("str").getValue()
    strStat.setBase(sourceStr)
    ctx.target.removeEffect(effectCtx.name)
end
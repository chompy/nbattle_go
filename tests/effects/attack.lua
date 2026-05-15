function Name()
    return "attack"
end

function OnAdd(effectCtx)
    effectCtx.target.getStat("hp").subBase(effectCtx.source.getStat("str").getValue() - effectCtx.target.getStat("def").getValue())
    effectCtx.remove()
end
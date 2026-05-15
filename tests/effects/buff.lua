function Name()
    return "buff"
end

function OnAdd(effectCtx)
    effectCtx.target.getStat("str").addBase(10)
    effectCtx.remove()
end
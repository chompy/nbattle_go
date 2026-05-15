function Name()
    return "regenerate"
end

function OnTick(effectCtx, evt)
    effectCtx.target.getStat("hp").addBase(3)
end

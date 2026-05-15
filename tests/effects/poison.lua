function Name()
    return "poison"
end

function OnTick(effectCtx, evt)
    effectCtx.target.getStat("hp").subBase(2)
end

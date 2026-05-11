function Name()
    return "poison"
end

function OnTick(ctx, evt)
    ctx.target.getStat("hp").subBase(2)
end

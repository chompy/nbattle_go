function Name()
    return "regenerate"
end

function OnTick(ctx, evt)
    ctx.target.getStat("hp").addBase(3)
end

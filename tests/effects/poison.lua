function Name()
    return "poison"
end

function OnTick(ctx, evt)
    ctx.target.getStat("hp").subtract(2)
end

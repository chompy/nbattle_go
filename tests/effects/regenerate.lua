function Name()
    return "regenerate"
end

function OnTick(ctx, evt)
    ctx.target.getStat("hp").add(3)
end

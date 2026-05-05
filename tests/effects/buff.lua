function Name()
    return "buff"
end

function OnAdd(ctx)
    ctx.target.getStat("str").add(10)
    --ctx.target.removeEffect(ctx.effect)
end
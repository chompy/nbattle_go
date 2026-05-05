function Name()
    return "attack"
end

function OnAdd(ctx)
    ctx.target.getStat("hp").subtract(ctx.source.getStat("str").get() - ctx.target.getStat("def").get())
    --ctx.target.removeEffect(ctx.effect)
end
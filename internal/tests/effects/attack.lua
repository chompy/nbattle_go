function Name()
    return "attack"
end

function OnAdd(ctx)
    ctx.target.getStat("hp").subBase(ctx.source.getStat("str").getValue() - ctx.target.getStat("def").getValue())
    --ctx.target.removeEffect(ctx.effect)
end
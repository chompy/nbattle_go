function Name()
    return "example_effect"
end

function OnAdd(ctx)
    ctx.target.getStat("hp").setBase(25)
end

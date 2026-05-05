function Name()
    return "example_effect"
end

function OnAdd(ctx)
    ctx.target.getStat("hp").set(25)
end

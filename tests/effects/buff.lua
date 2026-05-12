function Name()
    return "buff"
end

function OnAdd(ctx)
    ctx.target.getStat("str").addBase(10)
    ctx.remove()
end
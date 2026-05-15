function Name()
    return "example_effect"
end

function OnAdd(effectCtx)
    effectCtx.target.getStat("hp").setBase(25)
end

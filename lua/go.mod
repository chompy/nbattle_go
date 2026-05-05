module github.com/chompy/nbattle_go/lua

go 1.26.2

replace github.com/chompy/nbattle_go v0.0.0 => ../

replace github.com/chompy/nbattle_go/event v0.0.0 => ../event

require (
	github.com/chompy/nbattle_go v0.0.0-20260504195119-7ff16033d4e8
	github.com/chompy/nbattle_go/event v0.0.0-20260504195119-7ff16033d4e8
	github.com/rosbit/luago v0.5.2
)

require github.com/rosbit/go-embedding-utils v0.4.1 // indirect

module github.com/chompy/nbattle_go/tests

go 1.26.2

replace github.com/chompy/nbattle_go v0.0.0 => ../

replace github.com/chompy/nbattle_go/event v0.0.0 => ../event

replace github.com/chompy/nbattle_go/lua v0.0.0 => ../lua

require (
	github.com/chompy/nbattle_go v0.0.0
	github.com/chompy/nbattle_go/event v0.0.0
	github.com/chompy/nbattle_go/lua v0.0.0
)

require (
	github.com/rosbit/go-embedding-utils v0.4.1 // indirect
	github.com/rosbit/luago v0.5.2 // indirect
)

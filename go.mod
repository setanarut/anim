module github.com/setanarut/anim

go 1.24.0

require github.com/hajimehoshi/ebiten/v2 v2.8.6

require (
	github.com/ebitengine/gomobile v0.0.0-20250209143333-6071a2a2351c // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.2 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

retract v1.1.0 // SetStateAndReset() bug

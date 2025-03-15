module github.com/mrjrieke/nute

go 1.23.0

require (
	gioui.org v0.0.0-20220318070519-8833a6738a3b
	github.com/ftbe/dawg v0.0.0-20131228112149-aadae8139481
	github.com/g3n/engine v0.2.0
	golang.org/x/mobile v0.0.0-20231127183840-76ac6878050a
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/go-text/render v0.2.0 // indirect
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	github.com/jeandeaual/go-locale v0.0.0-20240223122105-ce5225dcaa49 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.4.0 // indirect
	github.com/rymdport/portal v0.2.6 // indirect
	golang.org/x/sync v0.11.0 // indirect
)

require golang.org/x/exp/shiny v0.0.0-20230817173708-d852ddb80c63 // indirect

require (
	fyne.io/systray v1.11.0 // indirect
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.6 // indirect
	github.com/benoitkugler/textlayout v0.0.10 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/fredbi/uri v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fyne-io/gl-js v0.0.0-20220119005834-d2da28d9ccfe // indirect
	github.com/fyne-io/glfw-js v0.0.0-20240101223322-6e1efdc71b7a // indirect
	github.com/fyne-io/image v0.0.0-20220602074514-4956b0afb3d2 // indirect
	github.com/gioui/uax v0.2.1-0.20220325163150-e3d987515a12 // indirect
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.4
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jsummers/gobmp v0.0.0-20151104160322-e2ba15ffa76e // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/yuin/goldmark v1.7.1 // indirect
	//golang.org/x/net v0.0.0-20220708220712-1185a9018129 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto v0.0.0-20220714211235-042d03aeabc9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	fyne.io/fyne/v2 v2.5.2
	golang.org/x/exp v0.0.0-20250215185904-eff6e970281f
	golang.org/x/image v0.20.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

//require github.com/go-gl/glfw/v3.3.2/glfw v0.0.0-20211213063430-748e38ca8aec

require (
	github.com/faiface/mainthread v0.0.0-20171120011319-8b78f0a41ae3
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a
	golang.org/x/net v0.25.0 // indirect
)

// Uncomment for local development
//replace fyne.io/fyne/v2 v2.5.2 => ../fyne // Use nute_integrate branch
replace fyne.io/fyne/v2 v2.5.2 => github.com/mrjrieke/fyne/v2 v2.5.2-1

// replace gioui.org v0.0.0-20220318070519-8833a6738a3b => ../../gio // Use mashup_v1 branch

//replace github.com/g3n/engine v0.2.0 => ../g3n/engine // Use mashup_v1 branch

// replace github.com/fyne-io/glfw-js v0.0.0-20220120001248-ee7290d23504 => ../../glfw-js // Use mashup_v1 branch

//replace fyne.io/fyne/v2 v2.1.3 => github.com/mrjrieke/fyne/v2 v2.1.3-6

replace gioui.org v0.0.0-20220318070519-8833a6738a3b => github.com/mrjrieke/gio v0.0.0-20220406132257-ec1380c11ef0

replace github.com/g3n/engine v0.2.0 => github.com/mrjrieke/engine v0.2.1-0.20230107141038-8bd28c2897c4

replace github.com/fyne-io/glfw-js v0.0.0-20220120001248-ee7290d23504 => github.com/mrjrieke/glfw-js v0.0.0-20220409154018-95a896685cdb

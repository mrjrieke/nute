module github.com/mrjrieke/nute

go 1.17

require (
	gioui.org v0.0.0-20220318070519-8833a6738a3b
	github.com/g3n/engine v0.2.0
	golang.org/x/mobile v0.0.0-20220307220422-55113b94f09c
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
)

require (
	fyne.io/systray v1.9.1-0.20220331100914-9177bf851614 // indirect
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.6 // indirect
	github.com/Kodeworks/golang-image-ico v0.0.0-20141118225523-73f0f4cfade9 // indirect
	github.com/benoitkugler/textlayout v0.0.10 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/faiface/glhf v0.0.0-20181018222622-82a6317ac380 // indirect
	github.com/faiface/mainthread v0.0.0-20171120011319-8b78f0a41ae3 // indirect
	github.com/fredbi/uri v0.0.0-20181227131451-3dcfdacbaaf3 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/fyne-io/gl-js v0.0.0-20220119005834-d2da28d9ccfe // indirect
	github.com/fyne-io/glfw-js v0.0.0-20220120001248-ee7290d23504 // indirect
	github.com/gioui/uax v0.2.1-0.20220325163150-e3d987515a12 // indirect
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-gl/mathgl v0.0.0-20190713194549-592312d8590a // indirect
	github.com/go-text/typesetting v0.0.0-20220112121102-58fe93c84506 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/goki/freetype v0.0.0-20181231101311-fa8a33aabaff // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20211219123610-ec9572f70e60 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/icexin/gocraft v0.0.0-20220710135235-50535c9c92c7 // indirect
	github.com/icexin/gocraft-server v0.0.0-20200316021447-c466fe50ae44 // indirect
	github.com/ojrac/opensimplex-go v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/srwiley/oksvg v0.0.0-20200311192757-870daf9aa564 // indirect
	github.com/srwiley/rasterx v0.0.0-20200120212402-85cb7272f5e9 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/yuin/goldmark v1.4.0 // indirect
	golang.org/x/net v0.0.0-20220708220712-1185a9018129 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220714211235-042d03aeabc9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	honnef.co/go/js/dom v0.0.0-20210725211120-f030747120f2 // indirect
)

require (
	fyne.io/fyne/v2 v2.1.3
	golang.org/x/exp v0.0.0-20210722180016-6781d3edade3 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/sys v0.0.0-20220712014510-0a85c31ab51e // indirect
)

require (
	github.com/go-gl/glfw v0.0.0-20200222043503-6f7a984d4dc4
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20211213063430-748e38ca8aec
)

// Uncomment for local development
// replace fyne.io/fyne/v2 v2.1.3 => ../../fyne // Use mashup_v1 branch

// replace gioui.org v0.0.0-20220318070519-8833a6738a3b => ../../gio // Use mashup_v1 branch

// replace github.com/g3n/engine v0.2.0 => ../../g3n/engine // Use mashup_v1 branch

// replace github.com/fyne-io/glfw-js v0.0.0-20220120001248-ee7290d23504 => ../../glfw-js // Use mashup_v1 branch

replace fyne.io/fyne/v2 v2.1.3 => github.com/mrjrieke/fyne/v2 v2.1.3-3

replace gioui.org v0.0.0-20220318070519-8833a6738a3b => github.com/mrjrieke/gio v0.0.0-20220406132257-ec1380c11ef0

//replace github.com/g3n/engine v0.2.0 => github.com/mrjrieke/engine v0.2.1-0.20220429130921-1cb788c5a9f8
replace github.com/g3n/engine v0.2.0 => github.com/mrjrieke/engine v0.2.1-0.20220804122658-e9fdaea14d50

replace github.com/fyne-io/glfw-js v0.0.0-20220120001248-ee7290d23504 => github.com/mrjrieke/glfw-js v0.0.0-20220409154018-95a896685cdb

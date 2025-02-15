package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	pixels []byte
	img    *ebiten.Image
}

func NewGame() *Game {
	// Create an empty pixel array (RGBA format)
	pixels := make([]byte, screenWidth*screenHeight*4)

	// Create an empty image
	img := ebiten.NewImage(screenWidth, screenHeight)

	return &Game{pixels: pixels, img: img}
}

// SetPixel modifies the pixel array at (x, y)
func (g *Game) SetPixel(x, y int, c color.Color) {
	if x < 0 || x >= screenWidth || y < 0 || y >= screenHeight {
		return
	}
	r, gVal, b, a := c.RGBA()
	i := (y*screenWidth + x) * 4
	g.pixels[i] = byte(r >> 8)
	g.pixels[i+1] = byte(gVal >> 8)
	g.pixels[i+2] = byte(b >> 8)
	g.pixels[i+3] = byte(a >> 8)
}

func (g *Game) Update() error {
	// Example: Draw a red pixel at (100, 100)
	g.SetPixel(100, 100, color.RGBA{255, 0, 0, 255})

	// Apply pixels to the image
	g.img.ReplacePixels(g.pixels)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the modified image onto the screen
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Pixel Rendering Example")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

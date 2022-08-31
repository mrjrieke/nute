//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"embed"
	"flag"
	"log"
	"os"

	"github.com/mrjrieke/nute/g3nd/g3nworld"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
	"github.com/mrjrieke/nute/g3nd/worldg3n/g3nrender"
	"github.com/mrjrieke/nute/mashupsdk"
)

var worldCompleteChan chan bool

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	headless := flag.Bool("headless", false, "Run headless")
	custos := flag.Bool("custos", false, "Run in guardian mode.")
	torusLayout := flag.Bool("toruslayout", false, "Use torus layout insead.")
	flag.Parse()
	worldLog, err := os.OpenFile("world.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	mashupRenderer := &g3nrender.MashupRenderer{}
	var torusRenderer *g3nrender.TorusRenderer
	if *torusLayout {
		torusRenderer = &g3nrender.TorusRenderer{GenericRenderer: g3nrender.GenericRenderer{RendererType: g3nrender.LAYOUT}, ActiveColor: &g3ndpalette.DARK_RED}
	} else {
		torusRenderer = &g3nrender.TorusRenderer{ActiveColor: &g3ndpalette.DARK_RED}
	}
	backgroundRenderer := &g3nrender.BackgroundRenderer{
		CollaboratingRenderer: torusRenderer}
	mashupRenderer.AddRenderer("Torus", backgroundRenderer.CollaboratingRenderer)
	mashupRenderer.AddRenderer("Background", backgroundRenderer)

	worldApp := g3nworld.NewWorldApp(*headless, mashupRenderer)

	worldApp.InitServer(*callerCreds, *insecure)

	if *headless {
		DetailedElements := []*mashupsdk.MashupDetailedElement{
			{
				Basisid:     -1,
				State:       &mashupsdk.MashupElementState{Id: -1, State: int64(mashupsdk.Mutable)},
				Name:        "{0}-Torus",
				Alias:       "It",
				Description: "",
				Renderer:    "Torus",
				Genre:       "Solid",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    []int64{-2, -3},
			},
			{
				Basisid:     -2,
				State:       &mashupsdk.MashupElementState{Id: -2, State: int64(mashupsdk.Mutable)},
				Name:        "{0}-AxialCircle",
				Alias:       "Inside",
				Description: "",
				Renderer:    "Torus",
				Genre:       "Space",
				Subgenre:    "Ento",
				Parentids:   []int64{-1},
				Childids:    []int64{-3, -4},
			},
			{
				Id:          -3,
				State:       &mashupsdk.MashupElementState{Id: -3, State: int64(mashupsdk.Mutable)},
				Name:        "Up-Side-Down",
				Alias:       "Up-Side-Down",
				Description: "",
				Genre:       "Attitude",
				Subgenre:    "180,0,0",
				Parentids:   nil,
				Childids:    nil,
			},
			{
				Id:          5,
				State:       &mashupsdk.MashupElementState{Id: 5, State: int64(mashupsdk.Init)},
				Name:        "ToriOne",
				Description: "Tori",
				Renderer:    "Torus",
				Genre:       "Collection",
				Subgenre:    "Torus",
				Parentids:   []int64{},
				Childids:    []int64{8, 9, 10},
			},
			{
				Id:            6,
				State:         &mashupsdk.MashupElementState{Id: 6, State: int64(mashupsdk.Init)},
				Name:          "BackgroundScene",
				Description:   "Background scene",
				Renderer:      "Background",
				Colabrenderer: "Torus",
				Genre:         "Collection",
				Subgenre:      "",
				Parentids:     []int64{},
				Childids:      []int64{7},
			},
			{
				Id:            7,
				State:         &mashupsdk.MashupElementState{Id: 7, State: int64(mashupsdk.Init)},
				Name:          "Outside",
				Alias:         "Outside",
				Description:   "",
				Renderer:      "Background",
				Colabrenderer: "Torus",
				Genre:         "Space",
				Subgenre:      "Exo",
				Parentids:     nil,
				Childids:      nil,
			},
			{
				Id:          8,
				State:       &mashupsdk.MashupElementState{Id: 8, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity-One",
				Description: "",
				Renderer:    "Torus",
				Genre:       "Abstract",
				Subgenre:    "",
				Parentids:   []int64{5},
				Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
			},
			{
				Id:          9,
				State:       &mashupsdk.MashupElementState{Id: 9, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity-Two",
				Description: "",
				Renderer:    "Torus",
				Genre:       "Abstract",
				Subgenre:    "",
				Parentids:   []int64{5},
				Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
			},
			{
				Id:          10,
				State:       &mashupsdk.MashupElementState{Id: 10, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity-Three",
				Description: "",
				Renderer:    "Torus",
				Genre:       "Abstract",
				Subgenre:    "",
				Parentids:   []int64{5},
				Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
			},
		}
		generatedElements, genErr := worldApp.MSdkApiHandler.UpsertMashupElements(
			&mashupsdk.MashupDetailedElementBundle{
				AuthToken:        "",
				DetailedElements: DetailedElements,
			})

		if genErr != nil {
			log.Fatalf(genErr.Error(), genErr)
		} else {
			generatedElements.DetailedElements[3].State.State = int64(mashupsdk.Clicked)

			elementStateBundle := mashupsdk.MashupElementStateBundle{
				AuthToken:     "",
				ElementStates: []*mashupsdk.MashupElementState{generatedElements.DetailedElements[3].State},
			}

			worldApp.MSdkApiHandler.UpsertMashupElementsState(&elementStateBundle)
		}
		go worldApp.MSdkApiHandler.OnResize(&mashupsdk.MashupDisplayHint{Width: 1600, Height: 800})

	}

	// Initialize the main window.
	go worldApp.InitMainWindow()

	<-worldCompleteChan
}

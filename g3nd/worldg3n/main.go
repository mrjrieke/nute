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
	flag.Parse()
	worldLog, err := os.OpenFile("world.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	mashupRenderer := &g3nrender.MashupRenderer{}
	mashupRenderer.AddRenderer("Torus", &g3nrender.TorusRenderer{})

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
				Genre:       "Collection",
				Subgenre:    "Torus",
				Parentids:   []int64{},
				Childids:    []int64{7, 8},
			},
			{
				Id:          6,
				State:       &mashupsdk.MashupElementState{Id: 6, State: int64(mashupsdk.Init)},
				Name:        "Outside",
				Alias:       "Outside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Exo",
				Parentids:   nil,
				Childids:    nil,
			},
			{
				Id:          7,
				State:       &mashupsdk.MashupElementState{Id: 7, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity-One",
				Description: "",
				Genre:       "Abstract",
				Subgenre:    "",
				Parentids:   []int64{5},
				Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
			},
			{
				Id:          8,
				State:       &mashupsdk.MashupElementState{Id: 8, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity-Two",
				Description: "",
				Genre:       "Abstract",
				Subgenre:    "",
				Parentids:   []int64{5},
				Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
			},
			{
				Id:          9,
				State:       &mashupsdk.MashupElementState{Id: 9, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity-Three",
				Description: "",
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

	}

	// Initialize the main window.
	go worldApp.InitMainWindow()

	<-worldCompleteChan
}

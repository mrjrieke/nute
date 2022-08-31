//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"embed"
	"flag"
	"log"
	"os"

	"github.com/mrjrieke/nute/examples/helloworld/hfhud/hfworld"
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
	worldLog, err := os.OpenFile("hfworld.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	detailedElements := []*mashupsdk.MashupDetailedElement{
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
			Childids:    []int64{-2, 4},
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
			Childids:    []int64{4},
		},
		{
			Id:          4,
			State:       &mashupsdk.MashupElementState{Id: 4, State: int64(mashupsdk.Mutable)},
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
			Alias:       "All",
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

	hfWorld := hfworld.NewHFWorldApp(*headless, detailedElements, nil)

	hfWorld.InitServer(*callerCreds, *insecure)

	// Initialize the main window.
	go hfWorld.InitMainWindow()

	<-worldCompleteChan
}

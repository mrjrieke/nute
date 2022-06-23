//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"flag"
	"log"
	"os"

	"tini.com/nute/g3nd/g3nworld"
	"tini.com/nute/mashupsdk"
)

var worldCompleteChan chan bool

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

	worldApp := g3nworld.NewWorldApp(*headless)

	worldApp.InitServer(*callerCreds, *insecure)

	if *headless {
		DetailedElements := []*mashupsdk.MashupDetailedElement{
			{
				Basisid:     -1,
				State:       &mashupsdk.MashupElementState{Id: -1, State: int64(mashupsdk.Mutable)},
				Name:        "{0}-Torus",
				Description: "",
				Genre:       "Solid",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    []int64{-2, -3},
			},
			{
				Basisid:     -2,
				State:       &mashupsdk.MashupElementState{Id: -2, State: int64(mashupsdk.Mutable)},
				Name:        "{0}-AxialCircle",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Ento",
				Parentids:   []int64{-1},
				Childids:    nil,
			},
			{
				Basisid:     -3,
				State:       &mashupsdk.MashupElementState{Id: -4, State: int64(mashupsdk.Mutable)},
				Name:        "{0}-SharedAttitude",
				Description: "",
				Genre:       "Attitude",
				Subgenre:    "180,0,0",
				Parentids:   []int64{-1},
				Childids:    nil,
			},
			{
				Id:          4,
				State:       &mashupsdk.MashupElementState{Id: 2, State: int64(mashupsdk.Init)},
				Name:        "ToriOne",
				Description: "Tori",
				Genre:       "",
				Subgenre:    "",
				Parentids:   []int64{},
				Childids:    []int64{6},
			},
			{
				Id:          5,
				State:       &mashupsdk.MashupElementState{Id: 2, State: int64(mashupsdk.Init)},
				Name:        "TorusEntity",
				Description: "",
				Genre:       "Abstract",
				Subgenre:    "",
				Parentids:   []int64{5},
				Childids:    []int64{-1}, // -1 -- generated and replaced by server since it is immutable.
			},
			{
				Id:          6,
				State:       &mashupsdk.MashupElementState{Id: 2, State: int64(mashupsdk.Init)},
				Name:        "Outside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Exo",
				Parentids:   nil,
				Childids:    nil,
			},
		}
		generatedElements, genErr := worldApp.MSdkApiHandler.UpsertMashupElements(
			&mashupsdk.MashupDetailedElementBundle{
				AuthToken:        "",
				DetailedElements: DetailedElements,
			})

		if genErr != nil {
			log.Fatal(genErr)
		} else {
			//fmt.Println("Got something: %v\n", generatedElements)
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

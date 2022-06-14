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
				Id:          1,
				State:       &mashupsdk.MashupElementState{Id: 1, State: int64(mashupsdk.Init)},
				Name:        "Inside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Ento",
				Parentids:   []int64{3},
				Childids:    nil,
			},
			{
				Id:          2,
				State:       &mashupsdk.MashupElementState{Id: 2, State: int64(mashupsdk.Init)},
				Name:        "Outside",
				Description: "",
				Genre:       "Space",
				Subgenre:    "Exo",
				Parentids:   nil,
				Childids:    nil,
			},
			{
				Id:          3,
				State:       &mashupsdk.MashupElementState{Id: 3, State: int64(mashupsdk.Init)},
				Name:        "torus",
				Description: "",
				Genre:       "Solid",
				Subgenre:    "Ento",
				Parentids:   nil,
				Childids:    []int64{1, 4},
			},
			{
				Id:          4,
				State:       &mashupsdk.MashupElementState{Id: 4, State: int64(mashupsdk.Init)},
				Name:        "Up-Side-Down",
				Description: "",
				Genre:       "Attitude",
				Subgenre:    "180,0,0",
				Parentids:   []int64{1, 3},
				Childids:    nil,
			},
		}

		_, _ = worldApp.MSdkApiHandler.UpsertMashupElements(
			&mashupsdk.MashupDetailedElementBundle{
				AuthToken:        "",
				DetailedElements: DetailedElements,
			})

	}

	// Initialize the main window.
	go worldApp.InitMainWindow()

	<-worldCompleteChan
}

//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"embed"
	"flag"
	"log"
	"os"

	"github.com/mrjrieke/nute-core/mashupsdk"
	"github.com/mrjrieke/nute/g3nd/data"
	"github.com/mrjrieke/nute/g3nd/g3nworld"
	g3ndpalette "github.com/mrjrieke/nute/g3nd/palette"
	"github.com/mrjrieke/nute/g3nd/worldg3n/g3nrender"
	"github.com/mrjrieke/nute/mashupsdk/client"
)

var worldCompleteChan chan bool

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("tls-skip-validation", false, "Skip server validation")
	custos := flag.Bool("custos", false, "Run in guardian mode.")
	headless := flag.Bool("headless", false, "Run headless")
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

	worldApp := g3nworld.NewWorldApp(*headless, *custos, mashupRenderer, nil)
	DetailedElements := []*mashupsdk.MashupDetailedElement{}

	if *custos {
		worldApp.MashupContext = client.BootstrapInit("custos", worldApp.MSdkApiHandler, nil, nil, insecure)

		libraryElementBundle, upsertErr := worldApp.MashupContext.Client.GetElements(
			worldApp.MashupContext,
			&mashupsdk.MashupEmpty{AuthToken: worldApp.GetAuthToken()},
		)
		if upsertErr != nil {
			log.Printf("G3n Element initialization failure: %s\n", upsertErr.Error())
		}

		DetailedElements = libraryElementBundle.DetailedElements
	} else if *headless {
		DetailedElements = data.GetExampleLibrary()
	} else {
		// Running in 'server' mode means mashup elements will be posted to this server.
		worldApp.InitServer(*callerCreds, *insecure, 0)
	}

	if *custos || *headless {
		//
		// Generate concrete elements from library elements.
		//
		generatedElementsBundle, genErr := worldApp.MSdkApiHandler.UpsertElements(
			&mashupsdk.MashupDetailedElementBundle{
				AuthToken:        "",
				DetailedElements: DetailedElements,
			})

		if !*headless {
			//
			// Upsert concrete elements to custos
			//
			_, custosUpsertErr := worldApp.MashupContext.Client.UpsertElements(
				worldApp.MashupContext,
				&mashupsdk.MashupDetailedElementBundle{
					AuthToken:        worldApp.GetAuthToken(),
					DetailedElements: generatedElementsBundle.DetailedElements,
				})

			if custosUpsertErr != nil {
				log.Fatalf(custosUpsertErr.Error(), custosUpsertErr)
			}

		}

		if genErr != nil {
			log.Fatalf(genErr.Error(), genErr)
		} else {
			//
			// Pick an initial element to 'click'
			//
			generatedElementsBundle.DetailedElements[3].State.State = int64(mashupsdk.Clicked)

			elementStateBundle := mashupsdk.MashupElementStateBundle{
				AuthToken:     "",
				ElementStates: []*mashupsdk.MashupElementState{generatedElementsBundle.DetailedElements[3].State},
			}

			worldApp.MSdkApiHandler.TweakStates(&elementStateBundle)
		}
		go worldApp.MSdkApiHandler.OnDisplayChange(&mashupsdk.MashupDisplayHint{Width: 1600, Height: 800})
	}

	// Initialize the main window.
	worldApp.InitMainWindow()

	<-worldCompleteChan
}

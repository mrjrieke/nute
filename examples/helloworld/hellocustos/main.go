//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"embed"
	"flag"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjrieke/nute/custos/custosworld"
	"github.com/mrjrieke/nute/g3nd/data"
	"github.com/mrjrieke/nute/mashupsdk"
)

var worldCompleteChan chan bool

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

func TorusParser(custosWorldApp *custosworld.CustosWorldApp, childId int64, concreteElement *mashupsdk.MashupDetailedElement) {
	child := custosWorldApp.MashupDetailedElementLibrary[childId]
	if child != nil && child.Alias != "" {
		log.Printf("TorusParser lookup on: %s\n", child.Alias)
		if fwb, fwbOk := custosWorldApp.FyneWidgetElements[child.Alias]; fwbOk {
			if fwb.MashupDetailedElement != nil && fwb.GuiComponent != nil {
				fwb.MashupDetailedElement.Copy(child)
				fwb.GuiComponent.(*container.TabItem).Text = child.Name
			}
		}
	}

	if child != nil && len(child.GetChildids()) > 0 {
		for _, cId := range child.GetChildids() {
			TorusParser(custosWorldApp, cId, concreteElement)
		}
	}
}

func OutsideClone(custosWorldApp *custosworld.CustosWorldApp, childId int64, concreteElement *mashupsdk.MashupDetailedElement) {
	custosWorldApp.FyneWidgetElements["Outside"].MashupDetailedElement.Copy(concreteElement)
}

func DetailMappedTabItemFyneComponent(custosWorldApp *custosworld.CustosWorldApp, id string) *container.TabItem {
	de := custosWorldApp.FyneWidgetElements[id].MashupDetailedElement

	tabLabel := widget.NewLabel(de.Description)
	tabLabel.Wrapping = fyne.TextWrapWord
	tabItem := container.NewTabItem(id, container.NewBorder(nil, nil, layout.NewSpacer(), nil, container.NewVBox(tabLabel, container.NewAdaptiveGrid(2,
		widget.NewButton("Show", func() {
			// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
			if len(custosWorldApp.ElementLoaderIndex) > 0 {
				mashupIndex := custosWorldApp.ElementLoaderIndex[custosWorldApp.FyneWidgetElements[de.Alias].GuiComponent.(*container.TabItem).Text]
				custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement = custosWorldApp.MashupDetailedElementLibrary[mashupIndex]

				custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Hidden, false)
				if custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement.Genre == "Collection" {
					custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
				}
				custosWorldApp.FyneWidgetElements[de.Alias].OnStatusChanged()
			}
		}), widget.NewButton("Hide", func() {
			if len(custosWorldApp.ElementLoaderIndex) > 0 {
				// Workaround... mashupdetailedelement points at wrong element sometimes, but shouldn't!
				mashupIndex := custosWorldApp.ElementLoaderIndex[custosWorldApp.FyneWidgetElements[de.Alias].GuiComponent.(*container.TabItem).Text]
				custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement = custosWorldApp.MashupDetailedElementLibrary[mashupIndex]

				custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Hidden, true)
				if custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement.Genre == "Collection" {
					custosWorldApp.FyneWidgetElements[de.Alias].MashupDetailedElement.ApplyState(mashupsdk.Recursive, true)
				}
				custosWorldApp.FyneWidgetElements[de.Alias].OnStatusChanged()
			}
		})))),
	)
	return tabItem
}

//go:embed gophericon.png
var gopherIcon embed.FS

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	headless := flag.Bool("headless", false, "Run headless")
	flag.Parse()
	worldLog, err := os.OpenFile("custos.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	detailedElements := data.GetExampleLibrary()

	custosWorld := custosworld.NewCustosWorldApp(*headless, detailedElements, nil)
	custosWorld.CustomTabItemRenderer["Torus"] = TorusParser

	custosWorld.Title = "Hello Custos"
	custosWorld.MainWindowSize = fyne.NewSize(800, 100)
	gopherIconBytes, _ := gopherIcon.ReadFile("gophericon.png")
	custosWorld.Icon = fyne.NewStaticResource("Gopher", gopherIconBytes)

	custosWorld.DetailMappedFyneComponent("Outside",
		"The magnetic field at any point outside the toroid is zero.",
		"Outside",
		DetailMappedTabItemFyneComponent)

	custosWorld.DetailMappedFyneComponent("Up-Side-Down",
		"Torus is up-side-down",
		"",
		DetailMappedTabItemFyneComponent)

	custosWorld.DetailMappedFyneComponent("All",
		"A group of torus or a tori.",
		"",
		DetailMappedTabItemFyneComponent)

	if !custosWorld.Headless {
		custosWorld.InitServer(*callerCreds, *insecure)
	}

	// Initialize the main window.
	custosWorld.InitMainWindow()

	<-worldCompleteChan
}

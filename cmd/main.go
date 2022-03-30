package main

import(
	"flag"
	"github.com/oliread/usbdmx"
	"github.com/oliread/usbdmx/ft232"
	"log"
	"time"
)

func main(){

	vid := uint16(0x0403)
	pid := uint16(0x6001)
	outputInterfaceID := flag.Int("output-id", 2, "Output interface ID for device")
	inputInterfaceID := flag.Int("input-id", 1, "Output interface ID for device")
	debugLevel := flag.Int("debug", 0, "Debug level for USB context")
	flag.Parse()

	// Create a configuration from our flags
	vid := uint16(0x0403)
	pid := uint16(0x6001)
	config := usbdmx.NewConfig(vid, pid, *outputInterfaceID,*inputInterfaceID, *debugLevel)

	// Get a usb context for our configuration
	config.GetUSBContext()

	// Create a controller and connect to it
	controller := ft232.NewDMXController(config)
	if err := controller.Connect(); err != nil {
		log.Fatalf("Failed to connect DMX Controller: %s", err)
	}

	go func(c *ft232.DMXController) {
		for {
			if err := controller.Render(); err != nil {
				log.Fatalf("Failed to render output: %s", err)
			}

			time.Sleep(30 * time.Millisecond)
		}
	}(&controller)


// set a fixture to bright white
controller.SetChannel(1, 0)
controller.SetChannel(2, 0)
controller.SetChannel(3, 0)

if err := controller.Render(); err != nil {
log.Fatal(err)
}


select{}

}
package main

import (
	"log"
	"net"
	"time"

	"github.com/collinux/gohue"
	"go.evanpurkhiser.com/netgear"

	"go.evanpurkhiser.com/aauto/modules/httplights"
	"go.evanpurkhiser.com/aauto/modules/lightson"
)

func main() {
	config := ReadConfig()

	hueBridge, err := hue.NewBridge(config.BridgeAddr)
	if err != nil {
		log.Printf("Unable to connect to hue bridge: %s\n", err)
		return
	}

	err = hueBridge.Login(config.BridgeLogin)
	if err != nil {
		log.Printf("Cannot login to hue bridge: %s\n", err)
		return
	}

	log.Printf("Hue Bridge authenticated.")

	// Configure WiFi phone connection light trigger
	netgearClient := netgear.NewClient(
		config.NetgearAddr,
		config.NetgearUser,
		config.NetgearPass,
	)

	phoneMAC, _ := net.ParseMAC(config.PhoneMAC)

	lightsonService := lightson.DeviceLightsTrigger{
		HueBridge:          hueBridge,
		NetgearClient:      netgearClient,
		SceneName:          "home",
		RouterPollInterval: time.Second * 30,
		DebouceInterval:    time.Minute * 2,
		TriggerDeviceMAC:   phoneMAC,
	}

	err = lightsonService.Start()
	if err != nil {
		log.Fatalf("Could not start listening: %q\n", err)
	}

	log.Printf("LightsOn WiFi MAC trigger started.")

	httpServer := httplights.Server{
		HueBridge:  hueBridge,
		ServerAddr: config.HTTPServerAddr,
	}

	// Configure desktop sleep trigger
	httpServer.RegisterModule(&httplights.DesktopTrigger{
		SceneName: "pre-sleep",

		DefaultSchedule: httplights.DesktopTriggerSchedule{
			StartTime: "11:00PM",
			ActiveFor: 8 * time.Hour,
		},
	})

	// Configure generic scene selector HTTP module
	httpServer.RegisterModule(&httplights.SelectScene{})

	httpServer.Start()

	log.Printf("Http server started.")

	<-make(chan int)
}

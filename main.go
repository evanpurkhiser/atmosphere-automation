package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/collinux/gohue"
	"go.evanpurkhiser.com/netgear"

	"go.evanpurkhiser.com/aauto/modules/httplights"
	"go.evanpurkhiser.com/aauto/modules/lightson"
)

var (
	bridgeAddr  = flag.String("bridge-addr", "", "Hue Bridge address")
	bridgeLogin = flag.String("bridge-login", "", "Hue Bridge login string")

	phoneMAC    = flag.String("phone-mac", "", "Phone MAC address for light triggering")
	netgearAddr = flag.String("netgear-addr", "", "Netgear router address for device connect light triggering")
	netgearUser = flag.String("netgear-user", "", "Netgear router username for device connect light triggering")
	netgearPass = flag.String("netgear-pass", "", "Netgear router password for device connect light triggering")
)

func main() {
	flag.Parse()

	hueBridge, err := hue.NewBridge(*bridgeAddr)
	if err != nil {
		log.Printf("Unable to connect to hue bridge: %s", err)
		return
	}

	err = hueBridge.Login(*bridgeLogin)
	if err != nil {
		log.Printf("Cannot login to hue bridge: %s", err)
		return
	}

	// Configure WiFi phone connection light trigger
	phoneMac, _ := net.ParseMAC(*phoneMAC)
	netgearClient := netgear.NewClient(*netgearAddr, *netgearPass, *netgearPass)

	lightsonService := lightson.DeviceLightsTrigger{
		HueBridge:          hueBridge,
		NetgearClient:      netgearClient,
		SceneName:          "home",
		RouterPollInterval: time.Second * 30,
		DebouceInterval:    time.Minute * 2,
		TriggerDeviceMAC:   phoneMac,
	}

	err = lightsonService.Start()
	if err != nil {
		log.Fatalf("Could not start listening: %q", err)
	}

	// Configure HTTP light control server
	httpServer := httplights.Server{
		HueBridge: hueBridge,
	}

	httpServer.RegisterModule(&httplights.SelectScene{})

	httpServer.RegisterModule(&httplights.DesktopTrigger{
		SceneName: "pre-sleep",

		DefaultSchedule: httplights.DesktopTriggerSchedule{
			StartTime: "11:00PM",
			ActiveFor: 8 * time.Hour,
		},
	})

	httpServer.Start()

	<-make(chan int)
}

package app

import (
	"net"
	"time"

	"github.com/collinux/gohue"
	"github.com/sirupsen/logrus"
	"go.evanpurkhiser.com/netgear"

	"go.evanpurkhiser.com/aauto/modules/httplights"
	"go.evanpurkhiser.com/aauto/modules/lightson"
)

func StartApp() {
	config := ReadConfig()

	hueBridge, err := hue.NewBridge(config.BridgeAddr)
	if err != nil {
		logrus.Fatalf("Unable to connect to hue bridge: %s\n", err)
	}

	err = hueBridge.Login(config.BridgeLogin)
	if err != nil {
		logrus.Fatalf("Cannot login to hue bridge: %s\n", err)
	}

	logrus.Infof("Hue Bridge authenticated.")

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
		logrus.Fatalf("Could not start listening: %s\n", err)
	}

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
}

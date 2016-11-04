package lightson

import (
	"net"
	"time"

	"github.com/collinux/gohue"
	"go.evanpurkhiser.com/netgear"
)

// DeviceLightsTrigger is a service that listens for a device to connect or
// disconnect from the network and will trigger a specified hue scene.
type DeviceLightsTrigger struct {
	// HueBridge specifies the bridge to communicate with for light changes.
	HueBridge *hue.Bridge

	// NetgearClient specifies the client to the router that will be used to
	// query for changes to the list of connected devices.
	NetgearClient *netgear.Client

	// TriggerDeviceMAC is the hardware address of the device that triggers the
	// service to turn the lights on or off.
	TriggerDeviceMAC net.HardwareAddr

	// SceneName specifies the name of the scene to trigger when the device
	// connects to the network.
	SceneName string

	// RouterPollInterval specifies the time between queries to the router to
	// determine if the device has been connected or disconnected.
	RouterPollInterval time.Duration

	// DebouceInterval specifies the time to wait before powering the lights
	// off. This allows the service to wait to ensure the device is not
	// reconnected to the network, as some devices tend to disconnect and
	// reconnect within a short period of time.
	DebouceInterval time.Duration
}

// lightsOff turns all lights off. This will wait wait before turning off the
// lights as it's presumed I won't be home to care, however this timer may be
// canceled should the lights be turned back on.
func (dt *DeviceLightsTrigger) lightsOff(cancel chan bool) {
	timer := time.NewTimer(dt.DebouceInterval)

	select {
	case <-cancel:
		timer.Stop()
		return
	case <-timer.C:
		break
	}

	nope := false
	dt.HueBridge.SetGroupState(0, &hue.Action{On: &nope})
}

// lightsOn sets the lights to the specified scene. This will only recall the
// scene given that all lights are currently off.
func (dt *DeviceLightsTrigger) lightsOn() {
	lights, _ := dt.HueBridge.GetAllLights()

	// Do nothing if any of the lights are currently on
	for _, light := range lights {
		if light.State.On {
			return
		}
	}

	dt.HueBridge.RecallSceneByName(dt.SceneName)
}

// Start boots the service and begins listening for devices to trigger lights.
func (dt *DeviceLightsTrigger) Start() error {
	// Ensure a valid scene was given
	if _, err := dt.HueBridge.GetSceneByName(dt.SceneName); err != nil {
		return err
	}

	cancelPowerOff := make(chan bool, 1)

	listener := func(change *netgear.ChangedDevice, err error) {
		if err != nil {
			return
		}

		if change.Device.MAC.String() != dt.TriggerDeviceMAC.String() {
			return
		}

		if change.Change == netgear.DeviceRemoved {
			go dt.lightsOff(cancelPowerOff)
			return
		}

		cancelPowerOff <- true
		close(cancelPowerOff)
		cancelPowerOff = make(chan bool, 1)

		go dt.lightsOn()
	}

	pollTime := time.Second * 10
	dt.NetgearClient.OnDeviceChanged(pollTime, listener)

	// TODO: Add a DeviceLightsTrigger.Stop() method which ensures the
	//       OnDeviceChanges call is also stopped. Currently this method
	//       returns a ticker, which can be stopped however it will leave a go
	//       routine in deadlock.

	return nil
}

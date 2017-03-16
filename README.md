## Atmosphere Automation service

This is a small service that I've written to automate various aspects of my
home in a programmatic way.

Current features:

 * Turn lights on / off when phone connects to the WiFi network.

 * Simple HTTP API for selecting light scenes.

 * API trigger for the 'computer standby' event. Determining if it's late
   enough to trigger the 'pre-sleep' light scene.

Usage:

The following environment variables must be configured:


```sh
# Hue Bridge address and login key
BRIDGE_ADDR=192.168.1.100
BRIDGE_LOGIN=someLoginKey

# MAC address of the device to listen for on the WiFi network
PHONE_MAC=C0:EE:FB:5C:44:28

# Netgear router authentication (for listening for the device)
NETGEAR_ADDR=192.168.1.1
NETGEAR_USER="admin"
NETGEAR_PASS="password"

# The HTTP API binding address
HTTP_SERVER_ADDR=:8080
```

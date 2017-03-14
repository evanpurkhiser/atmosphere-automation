package main

import (
	"github.com/kelseyhightower/envconfig"
)

// Config specifies the application configuration options.
type Config struct {
	// Hue bridge
	BridgeAddr  string `envconfig:"BRIDGE_ADDR"`
	BridgeLogin string `envconfig:"BRIDGE_LOGIN"`

	// Neatgear router authentication
	NetgearAddr string `envconfig:"NETGEAR_ADDR"`
	NetgearUser string `envconfig:"NETGEAR_USER"`
	NetgearPass string `envconfig:"NETGEAR_PASS"`

	// PhoneMAC address to listen for
	PhoneMAC string `envconfig:"PHONE_MAC"`

	// HTTPServerAddr specifies the address to bind to for the HTTP server
	HTTPServerAddr string `envconfig:"HTTP_SERVER_ADDR"`
}

// ReadConfig loads the configuration from the environment.
func ReadConfig() Config {
	config := Config{}
	envconfig.MustProcess("", &config)

	return config
}

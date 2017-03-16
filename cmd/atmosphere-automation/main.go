package main

import "go.evanpurkhiser.com/aauto/app"

func main() {
	app.StartApp()

	<-make(chan bool)
}

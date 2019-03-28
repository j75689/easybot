package main

import "easybot/server"

var (
	version string
	mode    = "develop"
)

func main() {
	server.Start(mode)
}

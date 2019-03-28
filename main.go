package main

import "github.com/j75689/easybot/server"

var (
	version string
	mode    = "develop"
)

func main() {
	server.Start(mode)
}

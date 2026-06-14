package main

import (
	_ "embed"

	_ "github.com/share-group/share-go/examples/daemon/crontab"

	"github.com/share-group/share-go/server"
)

//go:embed banner
var banner string

func main() {
	s := server.NewDaemonServer()
	s.SetBanner(banner)
	s.Run()
}

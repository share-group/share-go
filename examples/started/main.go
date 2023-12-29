package main

import (
	system "github.com/share-group/share-go/examples/started/system/controller"
	"github.com/share-group/share-go/server"
)

func main() {
	s := server.NewHttpServer()
	s.SetBanner("\n _______    _                                      _____            _     _                         _ \n|__   __|  | |                                    |  __ \\          | |   | |                       | |\n   | |_   _| |__   ___ _ __ ___   __ ___  ________| |  | | __ _ ___| |__ | |__   ___   __ _ _ __ __| |\n   | | | | | '_ \\ / _ \\ '_ ` _ \\ / _` \\ \\/ /______| |  | |/ _` / __| '_ \\| '_ \\ / _ \\ / _` | '__/ _` |\n   | | |_| | |_) |  __/ | | | | | (_| |>  <       | |__| | (_| \\__ \\ | | | |_) | (_) | (_| | | | (_| |\n   |_|\\__,_|_.__/ \\___|_| |_| |_|\\__,_/_/\\_\\      |_____/ \\__,_|___/_| |_|_.__/ \\___/ \\__,_|_|  \\__,_|")
	s.SetHandlers(system.SystemController)
	s.Run()
}

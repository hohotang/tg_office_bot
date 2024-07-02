package main

import (
	"log"
	"os"
	"tgbot/app"
)

var version string = "0"

func main() {
	tgApp := app.NewApp(version)
	tgApp.Run(true)
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

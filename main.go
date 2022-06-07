package main

import (
	"log"

	"io.twasyl/devcore/cmd"
	"io.twasyl/devcore/pkg/config"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Execute()
	err = config.Save()
	if err != nil {
		log.Fatal(err)
	}
}

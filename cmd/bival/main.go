package main

import (
	"log"
	"os"

	"github.com/neatflowcv/bival/internal/app/cli"
)

func main() {
	err := cli.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

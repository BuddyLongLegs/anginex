package main

import (
	"log"

	"github.com/BuddyLongLegs/anginex/src/proxy"
)

func main() {
	if err := proxy.Run(); err != nil {
		log.Fatal(err)
	}
}

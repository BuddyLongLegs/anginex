package main

import (
	"log"

	"github.com/BuddyLongLegs/anginex/src/config"
	"github.com/BuddyLongLegs/anginex/src/logger"
	"github.com/BuddyLongLegs/anginex/src/proxy"
)

func main() {
	conf := config.GetConfig()
	_ = conf

	dbLogger := &logger.DBLogger{}
	dbLogger.Conn()
	dbLogger.CreateTables()
	defer dbLogger.Close()

	if err := proxy.Run(conf, dbLogger); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/BuddyLongLegs/anginex/src/analytics"
	"github.com/BuddyLongLegs/anginex/src/config"
	"github.com/BuddyLongLegs/anginex/src/logger"
	"github.com/BuddyLongLegs/anginex/src/proxy"
)

func main() {
	conf := config.GetConfig()

	dbLogger := &logger.DBLogger{}
	dbLogger.Conn()
	dbLogger.InitDatabase()
	defer dbLogger.Close()

	go func() {
		analytics.AnalyticsAPI()
	}()

	if err := proxy.Run(conf, dbLogger); err != nil {
		log.Fatal(err)
	}
}

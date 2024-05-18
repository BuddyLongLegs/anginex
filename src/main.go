package main

import (
	"log"

	"github.com/BuddyLongLegs/anginex/src/analytics"
	"github.com/BuddyLongLegs/anginex/src/config"
	"github.com/BuddyLongLegs/anginex/src/logger"
	"github.com/BuddyLongLegs/anginex/src/osmetrics"
	"github.com/BuddyLongLegs/anginex/src/proxy"
)

func main() {
	conf := config.GetConfig()

	dbLogger := &logger.DBLogger{}
	dbLogger.Conn()
	dbLogger.InitDatabase()
	defer dbLogger.Close()

	// start the analytics server
	go func() {
		analytics.AnalyticsAPI(conf)
	}()

	// start the os metrics logger
	osMetrics := &osmetrics.OsMetric{}
	osMetrics.Init(!conf.Analytics.DisableSystemMetrics, dbLogger)
	go osMetrics.LogMetrics()

	if err := proxy.Run(conf, dbLogger); err != nil {
		log.Fatal(err)
	}
}

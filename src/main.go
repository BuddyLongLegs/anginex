package main

import (
	"log"

	"github.com/BuddyLongLegs/anginex/src/logger"
	"github.com/BuddyLongLegs/anginex/src/proxy"
)

func main() {
	dbLogger := &logger.DBLogger{}
	dbLogger.Conn()
	dbLogger.CreateTables()
	defer dbLogger.Close()

	addLog := dbLogger.AddLog

	if err := proxy.Run(addLog); err != nil {
		log.Fatal(err)
	}
}

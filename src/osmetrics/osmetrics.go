package osmetrics

import (
	"fmt"
	"time"

	"github.com/BuddyLongLegs/anginex/src/logger"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

type OsMetric struct {
	enabled   bool
	cpuBefore *cpu.Stats
	dbLogger  *logger.DBLogger
}

func (o *OsMetric) Init(enabled bool, dbLogger *logger.DBLogger) {
	o.enabled = true
	o.dbLogger = dbLogger
}

func (o *OsMetric) LogMetrics() {
	if !o.enabled {
		return
	}

	metricTime := time.Now().Local().Format(time.RFC3339)

	cpu, err := cpu.Get()
	if err != nil {
		fmt.Println("Error getting CPU stats:", err)
		return
	}

	mem, err := memory.Get()
	if err != nil {
		fmt.Println("Error getting memory stats:", err)
		return
	}

	memUsed := float64(mem.Used) / float64(mem.Total) * 100

	var cpuUsed float64 = 0

	if o.cpuBefore != nil {
		total := float64(cpu.Total - o.cpuBefore.Total)
		idle := float64(cpu.Idle - o.cpuBefore.Idle + cpu.Steal - o.cpuBefore.Steal)
		cpuUsed = 100 * (total - idle) / total
	}

	go func() {
		o.dbLogger.AddSystemMetrics(metricTime, cpuUsed, memUsed)
	}()

	o.cpuBefore = cpu
	// get float value of cpu used
	// wait for 1 sec
	time.Sleep(1 * time.Second)
	o.LogMetrics()
}

func (o *OsMetric) Disable() {
	o.enabled = false
}

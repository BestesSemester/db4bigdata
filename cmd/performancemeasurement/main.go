package main

import (
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/db"
	"git.sys-tem.org/caos/db4bigdata/internal/performancemeasurement"
	"git.sys-tem.org/caos/db4bigdata/internal/util"
)

func main() {
	util.SetupLogs()
	pm := performancemeasurement.New(db.MSQL, "horrorlog")
	pm.Start("test", 1*time.Second)
	//pm.MeasureRAM("testRAM", 1*time.Second)
	// pm.MeasureRAM("test")
	test(pm)
	time.Sleep(5 * time.Second)
	//pm.StopMeasureCPU()
	//pm.StopMeasureRAM()
	pm.Stop()
	pm.Run()
}

func test(pm performancemeasurement.PerformanceMeasurement) {
	defer pm.MeasureTime("test", time.Now())
	time.Sleep(5 * time.Second)
}

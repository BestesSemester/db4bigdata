package main

import (
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/performancemeasurement"
	"git.sys-tem.org/caos/db4bigdata/internal/util"
)

func main() {
	util.SetupLogs()
	pm := performancemeasurement.New(model.MSQL, "horrorlog")
	pm.MeasureCPU("test", 1*time.Second)
	// pm.MeasureRAM("test")
	test(pm)
	time.Sleep(20 * time.Second)
	pm.StopMeasureCPU()
	pm.Run()
}

func test(pm performancemeasurement.PerformanceMeasurement) {
	defer pm.MeasureTime(time.Now(), "test")
	time.Sleep(5 * time.Second)
}

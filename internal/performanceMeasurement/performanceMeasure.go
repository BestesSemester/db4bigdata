package performanceMeasurement

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type PerformanceMeasurement struct {
	DatabaseType string
	LogFilePath string
}

func New(databaseType string, logFilePath string) PerformanceMeasurement {
	// Creates a new instance of the perfomanceTaker.
	t := PerformanceMeasurement{DatabaseType: databaseType, LogFilePath: logFilePath}
	return t
}
// MeasureTime - Measures how long a given function took to execute.
func (p *PerformanceMeasurement) MeasureTime(givenTime time.Time,operation string) {
	elapsed := time.Since(givenTime)
	logrus.Printf("TIME: %s took %s\n", operation, elapsed)
	toWrite := fmt.Sprintf("It took %s Seconds to do the %s-operation",elapsed,operation)
	p.writeToFile(toWrite)
}
// MeasureCPU - Measures how much CPU power was needed to complete the operation.
func (p *PerformanceMeasurement) MeasureCPU(operation string) {
	percent, _ := cpu.Percent(0, true)
	toWrite := fmt.Sprintf("It took %.2f cpu power to do the %s-operation", percent[cpu.CPUser], operation)
	p.writeToFile(toWrite)
}


// writeToFile - Writes string to log.file
func (p *PerformanceMeasurement) writeToFile(content string) {
	logFile, err := os.OpenFile(p.LogFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		logrus.Println(err)
	} else {
		logrus.SetOutput(logFile)
		logrus.Println(content)
	}
	err = logFile.Close()
	if err != nil {
		logrus.Println(err)
	}
}
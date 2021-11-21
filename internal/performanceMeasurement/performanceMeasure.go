package performanceMeasurement

import (
"fmt"
"github.com/sirupsen/logrus"
"log"
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
// writeToFile - Writes string to log.file
func (p *PerformanceMeasurement) writeToFile(content string) {
	file, err := os.OpenFile(p.LogFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		logrus.Println(err)
	} else {
		logger := log.New(file, "prefix", log.LstdFlags)
		logger.Println(content)
		fmt.Println("Done")
	}
	err = file.Close()
	if err != nil {
		logrus.Println(err)
	}
}
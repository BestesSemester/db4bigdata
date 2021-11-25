package performancemeasurement

import (
	"fmt"
	"os"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

type PerformanceMeasurement struct {
	DatabaseType            model.StorageType
	LogFilePath             string
	startMeasureCPUChannel  chan string
	startMeasureTimeChannel chan TimeMeasurementParameters
	stopChannel             chan bool
	logChannel              chan string
	processes               int
}

type TimeMeasurementParameters struct {
	StartTime time.Time
	Operation string
}

func New(databaseType model.StorageType, logFilePath string) PerformanceMeasurement {
	// Creates a new instance of the perfomanceTaker.
	p := PerformanceMeasurement{DatabaseType: databaseType, LogFilePath: logFilePath}
	p.startMeasureTimeChannel = make(chan TimeMeasurementParameters)
	p.startMeasureCPUChannel = make(chan string)
	p.stopChannel = make(chan bool)
	p.logChannel = make(chan string)
	p.processes = 0
	p.startWatchers()
	return p
}

func (p *PerformanceMeasurement) startWatchers() {
	go p.startFileWriter()
	go p.ReadMeasureTime()
	go p.ReadMeasureCPU()
}

func (p *PerformanceMeasurement) MeasureTime(now time.Time, operation string) {
	p.processes++
	p.startMeasureTimeChannel <- TimeMeasurementParameters{
		StartTime: now,
		Operation: "test",
	}
}

func (p *PerformanceMeasurement) MeasureCPU(operation string) {
	p.processes++
	p.startMeasureCPUChannel <- operation
}

func (p *PerformanceMeasurement) Run() {
	for {
		time.Sleep(1 * time.Second)
	}
}

// MeasureTime - Measures how long a given function took to execute.
func (p *PerformanceMeasurement) ReadMeasureTime() {
	for {
		config, more := <-p.startMeasureTimeChannel
		if !more {
			p.stopChannel <- true
		}
		elapsed := time.Since(config.StartTime)
		logrus.Printf("TIME: %s took %s\n", config.Operation, elapsed)
		prtstr := fmt.Sprintf("It took %s Seconds to do the %s-operation", elapsed, config.Operation)
		p.writeToFile(prtstr)
	}
}

// MeasureCPU - Measures how much CPU power was needed to complete the operation.
func (p *PerformanceMeasurement) ReadMeasureCPU() {
	for {
		operation, more := <-p.startMeasureCPUChannel
		if !more {
			p.stopChannel <- true
		}
		percent, _ := cpu.Percent(0, true)
		logrus.Println("measuring cpu usage")
		prtstr := fmt.Sprintf("It took %.2f cpu power to do the %s-operation", percent[0], operation)
		p.writeToFile(prtstr)
		logrus.Println(prtstr)
	}
}

// writeToFile - Writes string to log.file
func (p *PerformanceMeasurement) writeToFile(content string) {
	p.processes++
	p.logChannel <- content
}

func (p *PerformanceMeasurement) startFileWriter() {
	logFile, err := os.OpenFile(p.LogFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		logrus.Println(err)
	}
	for {
		content, more := <-p.logChannel
		if !more {
			p.stopChannel <- true
		}
		if _, err := logFile.WriteString(content); err != nil {
			logrus.Errorln(err)
			break
		}
	}

	err = logFile.Close()
	if err != nil {
		logrus.Println(err)
	}
}

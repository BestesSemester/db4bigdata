package performancemeasurement

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

type PerformanceMeasurement struct {
	DatabaseType            model.StorageType
	LogFilePath             string
	startMeasureCPUChannel  chan TimeMeasurementParameters
	startMeasureRAMChannel  chan string
	startMeasureTimeChannel chan TimeMeasurementParameters
	stopChannel             chan bool
	stopChannelCPU          chan bool
	logChannel              chan string
}

type TimeMeasurementParameters struct {
	StartTime   time.Time
	Operation   string
	StopChannel chan bool
}

func New(databaseType model.StorageType, logFilePath string) PerformanceMeasurement {
	// Creates a new instance of the perfomanceTaker.
	p := PerformanceMeasurement{DatabaseType: databaseType, LogFilePath: logFilePath}
	p.startMeasureTimeChannel = make(chan TimeMeasurementParameters)
	p.startMeasureCPUChannel = make(chan TimeMeasurementParameters)
	p.startMeasureRAMChannel = make(chan string)
	p.stopChannel = make(chan bool)
	p.logChannel = make(chan string)
	p.startWatchers()
	return p
}

func (p *PerformanceMeasurement) startWatchers() {
	go p.startFileWriter()
	go p.ReadMeasureTime()
	go p.ReadMeasureRAM()
}

func (p *PerformanceMeasurement) MeasureTime(now time.Time, operation string) {
	p.startMeasureTimeChannel <- TimeMeasurementParameters{
		StartTime: now,
		Operation: operation,
	}
}

func (p *PerformanceMeasurement) MeasureCPU(operation string, interval time.Duration) {
	p.stopChannelCPU = make(chan bool)
	go p.readMeasureCPU(operation, interval)
}

func (p *PerformanceMeasurement) StopMeasureCPU() {
	p.stopChannelCPU <- true
}

func (p *PerformanceMeasurement) MeasureRAM(operation string) {
	p.startMeasureRAMChannel <- operation
}
func (p *PerformanceMeasurement) Run() {
	for {
		time.Sleep(1 * time.Second)
	}
}

// ReadMeasureTime - Measures how long a given function took to execute.
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
func (p *PerformanceMeasurement) ReadMeasureRAM() {
	for {
		operation, more := <-p.startMeasureRAMChannel
		if !more {
			p.stopChannel <- true
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		logrus.Println("measuring RAM usage")
		prtstr := fmt.Sprintf("Alloc = %v MiB for %s.", m.Alloc/1024/1024, operation)
		p.writeToFile(prtstr)
		logrus.Println(prtstr)
	}

}

// readMeasureCPU - Measures how much CPU power was needed to complete the operation.
func (p *PerformanceMeasurement) readMeasureCPU(operation string, interval time.Duration) {
	logrus.Println("starting cpu measurement")
	logrus.Println(p.stopChannelCPU)
	for {
		select {
		case <-p.stopChannelCPU:
			logrus.Println("stopping cpu measurement")
			return
		default:
			logrus.Println("taking default route")
		}
		percent, _ := cpu.Percent(0, true)
		logrus.Println("measuring cpu usage")
		prtstr := fmt.Sprintf("It took %.2f cpu power to do the %s-operation", percent[0], operation)
		p.writeToFile(prtstr)
		logrus.Println(prtstr)
		time.Sleep(interval)
	}
}

// writeToFile - Writes string to log.file
func (p *PerformanceMeasurement) writeToFile(content string) {
	p.logChannel <- content
}

func (p *PerformanceMeasurement) startFileWriter() {
	logFile, err := os.OpenFile(p.LogFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		logrus.Println(err)
	}
	for {
		content := <-p.logChannel
		if _, err := logFile.WriteString(content + "\n"); err != nil {
			logrus.Errorln(err)
			break
		}
	}

	err = logFile.Close()
	if err != nil {
		logrus.Println(err)
	}
}

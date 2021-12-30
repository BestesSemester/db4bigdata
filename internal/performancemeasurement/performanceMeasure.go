package performancemeasurement

import (
	"fmt"
	"os"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/db"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

type PerformanceMeasurement struct {
	DatabaseType            db.StorageType
	LogFilePath             string
	startMeasureCPUChannel  chan TimeMeasurementParameters
	startMeasureRAMChannel  chan string
	startMeasureTimeChannel chan TimeMeasurementParameters
	stopChannel             chan bool
	stopChannelCPU          chan bool
	stopChannelRAM          chan bool
	logChannel              chan string
}

type TimeMeasurementParameters struct {
	StartTime   time.Time
	Operation   string
	StopChannel chan bool
}

func New(databaseType db.StorageType, logFilePath string) PerformanceMeasurement {
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
}
func (p *PerformanceMeasurement) Start(operation string, interval time.Duration) {
	p.stopChannelCPU = make(chan bool)
	go p.readMeasureCPU(operation, interval)
	p.stopChannelRAM = make(chan bool)
	go p.ReadMeasureRAM(operation, interval)
}

func (p *PerformanceMeasurement) Stop() {
	p.stopChannelCPU <- true
	p.stopChannelRAM <- true

}
func (p *PerformanceMeasurement) MeasureTime(operation string, now time.Time) {
	p.startMeasureTimeChannel <- TimeMeasurementParameters{
		StartTime: now,
		Operation: operation,
	}
}

func (p *PerformanceMeasurement) StopMeasureCPU() {
	p.stopChannelCPU <- true
}

//func (p *PerformanceMeasurement) MeasureRAM(operation string, interval time.Duration) {
//	p.stopChannelRAM = make(chan bool)
//	go p.ReadMeasureRAM(operation, interval)
//}
func (p *PerformanceMeasurement) StopMeasureRAM() {
	p.stopChannelRAM <- true
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
		// logrus.Printf("TIME: %s took %s\n", config.Operation, elapsed)
		prtstr := fmt.Sprintf("TIME:,%s", elapsed)
		p.writeToFile(prtstr)
	}
}

// ReadMeasureRAM - Measures how much RAM the system uses.
func (p *PerformanceMeasurement) ReadMeasureRAM(operation string, interval time.Duration) {
	for {
		select {
		case <-p.stopChannelRAM:
			return
		default:
		}
		memory, err := memory.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		prtstr := fmt.Sprintf("RAM:,%v", memory.Used)
		p.writeToFile(prtstr)
		// logrus.Println(prtstr)
		time.Sleep(interval)
	}
}

// readMeasureCPU - Measures how much CPU power was needed to complete the operation.
func (p *PerformanceMeasurement) readMeasureCPU(operation string, interval time.Duration) {
	// logrus.Println(p.stopChannelCPU)
	for {
		select {
		case <-p.stopChannelCPU:
			return
		default:
		}
		percent, _ := cpu.Percent(0, true)
		prtstr := fmt.Sprintf("CPU:,%.2f", percent[0])
		p.writeToFile(prtstr)
		// logrus.Println(prtstr)
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

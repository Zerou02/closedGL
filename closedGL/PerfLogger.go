package closedGL

import (
	"time"
)

type TimeStruct struct {
	start time.Time
	dur   float64
}
type PerfLogger struct {
	logMap      map[string]*TimeStruct
	PrettyPrint bool
}

func NewLogger() PerfLogger {
	return PerfLogger{logMap: map[string]*TimeStruct{}, PrettyPrint: true}
}

func (this *PerfLogger) Start(name string) {
	var now = TimeStruct{
		start: time.Now(),
	}
	this.logMap[name] = &now
}

func (this *PerfLogger) End(name string) {
	var now = time.Now()
	this.logMap[name].dur = now.Sub(this.logMap[name].start).Seconds()
}

func (this *PerfLogger) Print() {
	println("-----")
	for k, v := range this.logMap {
		print(k + ": ")
		if this.PrettyPrint {
			PrintlnFloat(float32(v.dur))
		} else {
			println(v.dur)
		}
	}
	println("-----")
}

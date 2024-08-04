package closedGL

import (
	"time"
)

type TimeStruct struct {
	start time.Time
	dur   float64
}
type PerfLogger struct {
	logMap               map[string]*TimeStruct
	order                []string
	PrettyPrint, Enabled bool
}

func NewLogger() PerfLogger {
	return PerfLogger{logMap: map[string]*TimeStruct{}, PrettyPrint: true, Enabled: true, order: []string{}}
}

func (this *PerfLogger) Start(name string) {
	if !this.Enabled {
		return
	}
	var now = TimeStruct{
		start: time.Now(),
	}
	this.logMap[name] = &now
	if !ContainsString(this.order, name) {
		this.order = append(this.order, name)
	}
}

func (this *PerfLogger) End(name string) {
	if !this.Enabled {
		return
	}
	var now = time.Now()
	this.logMap[name].dur = now.Sub(this.logMap[name].start).Seconds()
}

func (this *PerfLogger) Print() {
	if !this.Enabled {
		return
	}
	println("-----")
	for _, k := range this.order {
		var v = this.logMap[k]
		print(k + ": ")
		if this.PrettyPrint {
			PrintlnFloat(float32(v.dur))
		} else {
			println(v.dur)
		}
	}
	println("-----")
}

package closedGL

import (
	"fmt"
	"time"
)

type Profiler struct {
	timeMap map[string]time.Time
	log     bool
}

func newProfiler() Profiler {
	return Profiler{timeMap: map[string]time.Time{}, log: true}
}
func (this *Profiler) startTime(name string) {
	this.timeMap[name] = time.Now()
}

func (this *Profiler) endTime(name string) {
	if !this.log {
		return
	}
	var end = time.Now()
	var dur = end.Sub(this.timeMap[name])
	fmt.Printf("%s:%f\n", name, dur.Seconds())
}

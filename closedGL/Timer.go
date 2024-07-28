package closedGL

type Timer struct {
	currS, targetDur float32
	repeat           bool
	finished         bool
}

func NewTimer(durationSec float32, repeat bool) Timer {
	return Timer{
		currS:     0,
		targetDur: durationSec,
		repeat:    repeat,
		finished:  false,
	}
}

func (this *Timer) Process(delta float32) {
	if this.repeat && this.IsTick() {
		this.currS = 0
	}
	this.currS += delta
}

func (this *Timer) IsTick() bool {
	var retVal = !this.finished && this.currS >= this.targetDur
	if retVal {
		this.finished = true
	}
	return retVal
}

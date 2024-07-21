package closedGL

type Timer struct {
	currS, targetDur float32
	repeat           bool
}

func newTimer(durationSec float32, repeat bool) Timer {
	return Timer{
		currS:     0,
		targetDur: durationSec,
		repeat:    repeat,
	}
}

func (this *Timer) process(delta float32) {
	if this.repeat && this.isTick() {
		this.currS = 0
	}
	this.currS += delta
}

func (this *Timer) isTick() bool {
	return this.currS >= this.targetDur
}

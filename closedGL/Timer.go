package closedGL

type Timer struct {
	currS, targetDur float32
	repeat           bool
}

func NewTimer(durationSec float32, repeat bool) Timer {
	return Timer{
		currS:     0,
		targetDur: durationSec,
		repeat:    repeat,
	}
}

func (this *Timer) Process(delta float32) {
	if this.repeat && this.IsTick() {
		this.currS = 0
	}
	this.currS += delta
}

func (this *Timer) IsTick() bool {
	return this.currS >= this.targetDur
}

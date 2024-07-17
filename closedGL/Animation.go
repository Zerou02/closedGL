package closedGL

type Animation struct {
	progress         float32
	start, end, curr float32
	durationSec      float32
	elapsedSec       float32
	repeat           bool
	circular         bool
}

func NewAnimation(start, end, duration float32, repeat bool, circular bool) Animation {
	return Animation{start: start, progress: 0, end: end, curr: start, durationSec: duration, elapsedSec: 0, repeat: repeat, circular: circular}
}

func (this *Animation) Restart(start, end, duration float32, repeat bool, circular bool) {
	*this = Animation{start: start, progress: 0, end: end, curr: start, durationSec: duration, elapsedSec: 0, repeat: repeat, circular: circular}
}

func (this *Animation) Process(delta float32) {
	if this.repeat && this.progress >= 1 {
		this.elapsedSec = 0
		this.progress = 0
	}

	this.elapsedSec += delta
	this.elapsedSec = Clamp(0, this.durationSec, this.elapsedSec)
	this.progress = this.elapsedSec / this.durationSec
	this.curr = Lerp(this.start, this.end, this.progress)

	if this.circular && this.progress >= 1 {
		var tmp = this.start
		this.start = this.end
		this.end = tmp
		this.elapsedSec = 0
		this.progress = 0

	}
}

func (this *Animation) GetValue() float32 {
	return this.curr
}

func (this *Animation) IsFinished() bool {
	return !this.repeat && this.progress >= 1
}

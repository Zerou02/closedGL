package closed_gl

type IAnimation interface {
	processAnim()
}
type Animation struct {
	value           *float32
	durationSeconds float32
	currentSecond   float32
	stepPerSecond   float32
	start           float32
	end             float32
	repeat          bool
	Stopped         bool
	Finished        bool
}

type ContinousAnimation struct {
}

func newAnimation(value *float32, durationSeconds, start, end float32, repeat bool) Animation {
	var anim = Animation{value: value, durationSeconds: durationSeconds, currentSecond: 0, stepPerSecond: 0, start: start, end: end, repeat: repeat, Stopped: false, Finished: false}
	anim.stepPerSecond = (anim.end - anim.start) / durationSeconds
	return anim
}
func (this *Animation) process(delta float32) {
	if this.Stopped {
		return
	}
	if this.currentSecond >= this.durationSeconds {
		if this.repeat {
			*this.value = this.start
			this.currentSecond = 0
		} else {
			this.Finished = true
			return
		}
	}
	*this.value += delta * this.stepPerSecond
	this.currentSecond += delta

}

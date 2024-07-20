package closedGL

type StaggeredAnimation struct {
	anims    []Animation
	CurrIdx  int
	Finished bool
}

func NewStaggeredAnimation(anims []Animation) StaggeredAnimation {
	return StaggeredAnimation{
		anims: anims, CurrIdx: 0, Finished: false,
	}
}

func (this *StaggeredAnimation) Process(delta float32) {
	if this.Finished {
		return
	}
	var currAnim = &this.anims[this.CurrIdx]
	currAnim.Process(delta)
	if currAnim.IsFinished() {
		this.CurrIdx++
		if this.CurrIdx >= len(this.anims) {
			this.Finished = true
		}
	}
}

func (this *StaggeredAnimation) GetValue() float32 {
	if this.Finished {
		return 0
	}
	return this.anims[this.CurrIdx].GetValue()
}

func (this *StaggeredAnimation) GetValueArr() []float32 {
	var retArr = []float32{}
	for _, x := range this.anims {
		retArr = append(retArr, x.GetValue())
	}
	return retArr
}

func (this *StaggeredAnimation) IsFinished() bool {
	return this.Finished
}

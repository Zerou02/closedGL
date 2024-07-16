package closedGL

type State struct {
	TransitionConditions []func() bool
	TransitionFn         []func()
	TransitionIndexArray []int
	Function             func()
}

type Fsm struct {
	currentState int
	States       []State
	DrawFn       func()
}

func newFsm() Fsm {

	return Fsm{States: []State{}, currentState: 0, DrawFn: func() {}}
}

func newState(fn func()) State {
	return State{Function: fn, TransitionConditions: []func() bool{}, TransitionIndexArray: []int{}}
}

func (this *State) addTransition(fn func() bool, to int, transFn func()) {
	this.TransitionConditions = append(this.TransitionConditions, fn)
	this.TransitionIndexArray = append(this.TransitionIndexArray, to)
	this.TransitionFn = append(this.TransitionFn, transFn)
}

func (this *Fsm) Process() {
	if this.currentState >= len(this.States) {
		return
	}
	var currState = this.States[this.currentState]
	currState.Function()
	for i, x := range currState.TransitionConditions {
		if x() {
			currState.TransitionFn[i]()
			this.currentState = currState.TransitionIndexArray[i]
			break
		}
	}
}

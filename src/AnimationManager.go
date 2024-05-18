package main

type AnimationManager struct {
	anims map[string]*Animation
}

func newAnimationManger() AnimationManager {
	return AnimationManager{anims: map[string]*Animation{}}
}

func (this *AnimationManager) addAnim(id string, value *float32, duration, start, end float32, repeat bool) {
	var a = newAnimation(value, duration, start, end, repeat)
	this.anims[id] = &a
}

func (this *AnimationManager) cancelAnim(id string) {
	this.anims[id] = nil
}
func (this *AnimationManager) process(delta float32) {
	for k, v := range this.anims {
		if v == nil {
			continue
		}
		v.process(delta)
		if v.Finished {
			this.anims[k] = nil
		}
	}
}

package closedGL

/*

import "github.com/go-gl/glfw/v3.2/glfw"

type keyboardFun func()

type KeyInfo struct {
	pressed bool
}

type KeyBoardManager struct {
	window             *glfw.Window
	registeredKeysDown map[glfw.Key]bool
	currKeyDown        glfw.Key
}

func newKeyBoardManager(window *glfw.Window) KeyBoardManager {
	var manager = KeyBoardManager{window: window, registeredKeysDown: map[glfw.Key]bool{}}

	return manager
}

func (this *KeyBoardManager) IsDown(key glfw.Key) bool {
	return this.window.GetKey(key) == glfw.Press
}

func (this *KeyBoardManager) registerKey(key glfw.Key) {
	var isDown = this.IsDown(key)
	this.registeredKeysDown[key] = isDown

}
func (this *KeyBoardManager) Process() {
	if this.currKeyDown != 0 {
		this.registeredKeysDown[this.currKeyDown] = this.IsDown(this.currKeyDown)
		this.currKeyDown = 0
	}
	for x, y := range this.registeredKeysDown {
		var isDown = this.IsDown(x)
		if isDown && !y {
			this.currKeyDown = x
		} else {
			this.registeredKeysDown[x] = isDown
		}
	}
}

func (this *KeyBoardManager) IsPressed(key glfw.Key) bool {
	var _, isInMap = this.registeredKeysDown[key]
	if !isInMap {
		this.registeredKeysDown[key] = this.IsDown(key)
		return false
	}
	return this.currKeyDown == key
}
*/

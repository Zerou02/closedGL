package main

import "github.com/go-gl/glfw/v3.2/glfw"

type KeyBoardManager struct {
	window *glfw.Window
}

func newKeyBoardManager(window *glfw.Window) KeyBoardManager {
	var manager = KeyBoardManager{window: window}
	return manager
}

func (this *KeyBoardManager) isDown(key glfw.Key) bool {
	return this.window.GetKey(key) == glfw.Action(glfw.KeyDown)
}

func (this *KeyBoardManager) process() {

}

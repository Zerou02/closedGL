package closedGL

import "github.com/go-gl/gl/v4.1-core/gl"



type TextureMane struct {
	textures  []uint32
	handleMap map[string]uint64
}

func newTextureMane() TextureMane {
	return TextureMane{textures: []uint32{}, handleMap: map[string]uint64{}}
}

func (this *TextureMane) loadTex(path string) {
	if this.handleMap[path] != 0 {
		return
	}

	var texture = *LoadImage(path, gl.RGBA)
	var handle = gl.GetTextureHandleARB(texture)
	if handle == 0 {
		panic("invalid textureHandle")
	}
	this.textures = append(this.textures, texture)
	this.handleMap[path] = handle
}

func (this *TextureMane) getHandle(path string) uint64 {
	return this.handleMap[path]
}
func (this *TextureMane) makeResident() {
	for _, v := range this.handleMap {
		gl.MakeTextureHandleResidentARB(v)
	}
}

func (this *TextureMane) makeNonResident() {
	for _, v := range this.handleMap {
		gl.MakeTextureHandleNonResidentARB(v)
	}
}

func (this *TextureMane) copy() TextureMane {
	var newArr = make([]uint32, len(this.textures))
	copy(newArr, this.textures)
	var newMap = map[string]uint64{}
	for k, v := range this.handleMap {
		newMap[k] = v
	}
	return TextureMane{
		textures:  newArr,
		handleMap: newMap,
	}
}

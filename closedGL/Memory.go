package closedGL

import "github.com/EngoEngine/glm"

func extendArray(arr *[]float32, lenWithNewEntriesFloats int) {
	if lenWithNewEntriesFloats != 0 && len(*arr) == 0 {
		*arr = make([]float32, 1)
	}
	for lenWithNewEntriesFloats >= len(*arr) {
		var newArr = make([]float32, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

func extendArrayU16(arr *[]uint16, newLenEntries int) {

	if newLenEntries != 0 && len(*arr) == 0 {
		*arr = make([]uint16, 1)
	}
	for newLenEntries >= len(*arr) {
		var newArr = make([]uint16, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

func extendArrayVec4(arr *[]glm.Vec4, newLenEntries int) {
	if newLenEntries != 0 && len(*arr) == 0 {
		*arr = make([]glm.Vec4, 1)
	}
	for newLenEntries >= len(*arr) {
		var newArr = make([]glm.Vec4, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

func extendArrayU64(arr *[]uint64, newLenEntries int) {
	if newLenEntries != 0 && len(*arr) == 0 {
		*arr = make([]uint64, 1)
	}
	for newLenEntries >= len(*arr) {
		var newArr = make([]uint64, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

func extendArrayU32(arr *[]uint32, newLenEntries int) {
	if newLenEntries != 0 && len(*arr) == 0 {
		*arr = make([]uint32, 1)
	}
	for newLenEntries >= len(*arr) {
		var newArr = make([]uint32, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

func extendArrayU8(arr *[]uint8, newLenEntries int) {
	if newLenEntries != 0 && len(*arr) == 0 {
		*arr = make([]uint8, 1)
	}
	for newLenEntries >= len(*arr) {
		var newArr = make([]uint8, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

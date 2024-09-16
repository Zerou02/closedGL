package closedGL

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/EngoEngine/glm"
)

func printlnTexData(texData []byte) {
	for i := 0; i < 16*32*4*2; i += 4 {
		if texData[i+1] == 0xFF {
			print("1")
		} else {
			print("0")
		}
		if i%128 == 0 {
			println()
		}
	}
}

func printFloatArr(arr *[]float32, stride int) {
	for i := 0; i < len(*arr); i++ {
		if i%stride == 0 {
			println()
		}
		PrintFloat((*arr)[i])
		print(", ")
	}
	println()
}
func PrintFloat(f float32) {
	fmt.Printf("%f", f)
}

func PrintlnFloat(f float32) {
	fmt.Printf("%f\n", f)
}

func PrintlnVec2(vec glm.Vec2) {
	PrintFloat(vec[0])
	print(",")
	PrintlnFloat(vec[1])
}

func PrintByteArr(arr []byte, stride int) {
	for i := 0; i < len(arr); i++ {
		if i%stride == 0 {
			println()
		}
		fmt.Printf("%x ", arr[i])
	}
	println()
}

func RleEncode(arr []byte) []byte {
	var retBytes = []byte{}
	var count byte = 0
	var currByte = arr[0]
	for _, x := range arr {
		if x == currByte {
			count++
			if count == 0xff {
				retBytes = append(retBytes, count)
				retBytes = append(retBytes, currByte)
				count = 0
			}
		} else {
			retBytes = append(retBytes, count)
			retBytes = append(retBytes, currByte)
			currByte = x
			count = 1
		}
	}
	if count != 0 {
		retBytes = append(retBytes, count)
		retBytes = append(retBytes, currByte)
	}
	return retBytes
}

func RleDecode(arr []byte) []byte {
	var retBytes = []byte{}
	for i := 0; i < len(arr); i += 2 {
		var c = arr[i]
		var ch = arr[i+1]
		for j := 0; j < int(c); j++ {
			retBytes = append(retBytes, ch)
		}
	}
	return retBytes
}

func parseConfig(path string) map[string]string {
	var retMap = map[string]string{}
	var bytes, err = os.ReadFile(path)
	if err != nil {
		println("could not find config", err.Error())
	}
	var content = string(bytes)
	var linSep = "\n"
	if runtime.GOOS == "windows" {
		linSep = "\rn"
	}
	for _, x := range strings.Split(content, linSep) {
		var l = strings.Trim(x, " ")
		if l == "" || l[0] == '[' || l == " " {
			continue
		}
		var splitted = strings.Split(l, "=")
		retMap[splitted[0]] = splitted[1]
	}
	return retMap
}

func strToBool(str string) bool {
	if str == "true" {
		return true
	} else {
		return false
	}
}

func ContainsString(arr []string, x string) bool {
	for _, y := range arr {
		if y == x {
			return true
		}
	}
	return false
}

func Contains[T comparable](arr *[]T, x T) bool {
	var retVal = false
	for _, y := range *arr {
		if y == x {
			retVal = true
			break
		}
	}
	return retVal
}

func Remove[T comparable](arr []T, x T) []T {
	var retVal = []T{}
	for _, y := range arr {
		if y != x {
			retVal = append(retVal, y)
		}
	}
	return retVal
}

func FindIdx[T comparable](arr []T, x T) int {
	var retVal = -1
	for i, y := range arr {
		if x == y {
			retVal = i
			break
		}
	}
	return retVal
}

func FindLastIdx[T comparable](arr []T, x T) int {
	var retVal = -1
	for i, y := range arr {
		if x == y {
			retVal = i
		}
	}
	return retVal
}

func FindAmount[T comparable](arr []T, x T) int {
	var retVal = 0
	for _, y := range arr {
		if x == y {
			retVal++
		}
	}
	return retVal
}

func RemoveAt[T comparable](arr []T, idx int) []T {
	var retVal = []T{}
	for i, y := range arr {
		if i != idx {
			retVal = append(retVal, y)
		}
	}
	return retVal
}

// element becomes new one at idx
func InsertAt[T comparable](arr []T, x T, idx int) []T {
	if idx >= len(arr) {
		return append(arr, x)
	}
	var new = []T{}
	for i, y := range arr {
		if i == idx {
			new = append(new, x)
		}
		new = append(new, y)
	}
	return new
}

func InsertArrAt[T comparable](arr []T, x []T, idx int) []T {
	if idx >= len(arr) {
		return append(arr, x...)
	}
	var new = []T{}
	for i, y := range arr {
		if i == idx {
			for _, z := range x {
				new = append(new, z)
			}
		}
		new = append(new, y)
	}
	return new
}

func Reverse[T any](arr []T) []T {
	var retVal = []T{}
	for i := len(arr) - 1; i >= 0; i-- {
		retVal = append(retVal, arr[i])
	}
	return retVal
}

func Ternary[T any](cond bool, pos, neg T) T {
	if cond {
		return pos
	} else {
		return neg
	}
}

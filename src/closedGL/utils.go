package closed_gl

import "fmt"

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
	fmt.Printf("%f\n", f)
}

func PrintlnFloat(f float32) {
	fmt.Printf("%f\n", f)
}

func PrintlnVec2(vec Vec2) {
	PrintlnFloat(vec[0])
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

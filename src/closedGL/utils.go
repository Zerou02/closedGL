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
		printFloat((*arr)[i])
		print(", ")
	}
	println()
}
func printFloat(f float32) {
	fmt.Printf("%f\n", f)
}

func printlnFloat(f float32) {
	fmt.Printf("%f\n", f)
}

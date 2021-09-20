package main

import (
	"fmt"
)

var b64 = `src="image/based64TWFuIGlzIGRpc3RpAAA"`

func main() {
	// fmt.Print("test")
	// var validID = regexp.MustCompile(`image\/based64([^\"]*)(AAA)`)
	// result := validID.FindAllStringSubmatch(b64, -1)
	// fmt.Println(result[0][2])
	var test uint8 = 253
	for i := 0; i < 10; i++ {
		fmt.Println(int(test))
		test++
	}

}

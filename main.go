package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(strconv.IntSize)
	var number int64
	number = 99999
	fmt.Println(strconv.FormatInt(number, 2))
}














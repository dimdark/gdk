package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	msg := "Hello, 世界"
	bs := []byte(msg)
	for _, b := range bs {
		fmt.Printf("%x ", b)
	}
	fmt.Println()
	base64Msg := base64.StdEncoding.EncodeToString(bs)
	fmt.Println(base64Msg)
	bs, err := base64.StdEncoding.DecodeString(base64Msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bs))

	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	_, err = encoder.Write(bs)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(base64Msg)))
	var anobs [100]byte
	_, err = decoder.Read(anobs[:])
	fmt.Println(string(anobs[:]))
}














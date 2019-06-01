package base64

import (
	"encoding/base64"
	"fmt"
	"os"
)

func Example() {
	msg := "Hello, 世界"
	// base64编码
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	// base64解码
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(string(decoded))
	// Output:
	// SGVsbG8sIOS4lueVjA==
	// Hello, 世界
}

// ExampleStructName_MethodName
func ExampleEncoding_EncodeToString() {
	data := []byte("any + old & data")
	str := base64.StdEncoding.EncodeToString(data)
	fmt.Println(str)
	// Output:
	// YW55ICsgb2xkICYgZGF0YQ==
}

func ExampleEncoding_DecodeString() {
	str := "c29tZSBkYXRhIHdpdGggACBhbmQg77u/"
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%q\n", data)
	// Output:
	// "some data with \x00 and \ufeff"
}

// ExampleFunctionName
func ExampleNewEncoder() {
	input := []byte("foo\x00bar")
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	_, _ = encoder.Write(input)
	// 注意及时关闭
	_ = encoder.Close()
	// Output:
	// Zm9vAGJhcg==
}
















package driver

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// 用于测试ValueConverter接口及其实现的结构体
type valueConverterTest struct {
	c ValueConverter
	in interface{}
	out interface{}
	err string
}

var now = time.Now()
var answer int64 = 42

type (
	i int64
	f float64
	b bool
	bs []byte
	s string
	t time.Time
	is []int
)

func TestValueConverters(t *testing.T) {
	vct := valueConverterTest{c: Bool, in: "foo", err: "sql/driver: couldn't convert \"foo\" into type bool"}
	vct = valueConverterTest{c:Bool, in:"false", out:false, err:""}
	vct = valueConverterTest{DefaultParameterConverter, now, now, ""}
	vct = valueConverterTest{DefaultParameterConverter, (*int64)(nil), nil, ""}
	vct = valueConverterTest{DefaultParameterConverter, &answer, answer, ""}
	vct = valueConverterTest{DefaultParameterConverter, &now, now, ""}
	vct = valueConverterTest{DefaultParameterConverter, i(9), int64(9), ""}
	vct = valueConverterTest{DefaultParameterConverter, f(0.1), float64(0.1), ""}
	vct = valueConverterTest{DefaultParameterConverter, b(true), true, ""}
	vct = valueConverterTest{DefaultParameterConverter, bs{1}, []byte{1}, ""}
	vct = valueConverterTest{DefaultParameterConverter, s("a"), "a", ""}
	vct = valueConverterTest{DefaultParameterConverter, is{1}, nil, "unsupported type driver.is, a slice of int" }
	out, err := vct.c.ConvertValue(vct.in)
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	if errStr != vct.err {
		t.Errorf("%T(%T(%v)) error = %q; want error = %q", vct.c, vct.in, vct.in, errStr, vct.err)
	}
	if vct.err != "" {
		return
	}
	if !reflect.DeepEqual(out, vct.out) {
		t.Errorf("%T(%T(%v)) = %v (%T); want %v (%T)", vct.c, vct.in, vct.in, out, out, vct.out, vct.out)
	}
	fmt.Printf("%T(%T(%v)) = %v (%T); want %v (%T)", vct.c, vct.in, vct.in, out, out, vct.out, vct.out)
}





















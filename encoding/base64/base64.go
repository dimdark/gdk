package base64

import (
	"io"
	"strconv"
)

// base64编解码
type Encoding struct {
	// 编解码的字母表(映射表)
	encode [64]byte
	decodeMap [256]byte
	// 尾部填充的字符
	padChar rune
	strict bool
}

// 替换填充的字符
func (enc Encoding) WithPadding(padding rune) *Encoding {
	if padding == '\r' || padding == '\n' || padding > 0xff {
		panic("invalid padding")
	}
	for _, b := range enc.encode {
		if rune(b) == padding {
			panic("padding contained in alphabet")
		}
	}
	enc.padChar = padding
	return &enc
}

// 开启strict模式
func (enc Encoding) Strict() *Encoding {
	enc.strict = true
	return &enc
}

// base64编码
func (enc *Encoding) Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}
	di, si := 0, 0
	// 以3个字节为一个基本单位
	// 总共有n个基本单位
	n := (len(src) / 3) * 3
	for si < n {
		val := uint(src[si+0] << 16) | uint(src[si+1] << 8) | uint(src[si+2])
		dst[di+0] = byte(val >> 18)
		dst[di+1] = byte(val >> 12)
		dst[di+2] = byte(val >> 6)
		dst[di+3] = byte(val)
		si += 3
		di += 4
	}
	remain := len(src) - n
	if remain == 0 {
		return
	}
}


type encoder struct {
	err error
	enc *Encoding
	w io.Writer
	buf [3]byte
	nbuf int
	out [1024]byte
}

type decoder struct {
	err error
	readErr error
	enc *Encoding
	r io.Reader
	buf [1024]byte
	nbuf int
	out []byte
	outbuf [1024 / 4 * 3]byte
}

type CorruptInputError int64
func (e CorruptInputError) Error() string {
	return "illegal base64 data at input byte" + strconv.FormatInt(int64(e), 10)
}

type newlineFilteringReader struct {
	wrapped io.Reader
}

func (r *newlineFilteringReader) Read(p []byte) (int, error) {
	n, err := r.wrapped.Read(p)
	for n > 0 {
		offset := 0
		for i, b := range p {
			if b != '\r' && b != '\n' {
				if i != offset {
					p[offset] = b
				}
				offset++
			}
		}
		if offset > 0 {
			return offset, err
		}
		n, err = r.wrapped.Read(p)
	}
	return n, err
}

const (
	StdPadding rune = '='
	NoPadding rune = -1
)

const (
	encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	encodeUrl = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

// 标准 base64编解码器
var StdEncoding = NewEncoding(encodeStd)
// url base64编解码器
var UrlEncoding = NewEncoding(encodeUrl)
// 标准 base64编解码器 不用填充
var RawStdEncoding = StdPadding.WithPadding(NoPadding)
// url base64编解码器 不用填充
var RawUrlEncoding = UrlEncoding.WithPadding(NoPadding)

// 根据字母表(映射表)构建base64的Encoding
func NewEncoding(encoder string) *Encoding {
	if len(encoder) != 64 {
		panic("encoding alphabet is not 64-bytes long")
	}
	for i := 0; i < len(encoder); i++ {
		if encoder[i] == '\r' || encoder[i] == '\n' {
			panic("encoding alphabet contains newline character")
		}
	}

	e := new(Encoding)
	e.padChar = StdPadding
	// 切片复制
	copy(e.encode[:], encoder)
	for i := 0; i < len(e.decodeMap); i++ {
		// 0xFF -1
		e.decodeMap[i] = 0xFF
	}
	for i := 0; i < len(encoder); i++ {
		e.decodeMap[encoder[i]] = byte(i)
	}

	return e
}

func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	return &encoder{enc: enc, w: w}
}

func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	return &decoder{enc: enc, r: &newlineFilteringReader{r}}
}















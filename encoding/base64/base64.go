package base64

import (
	"encoding/binary"
	"io"
	"strconv"
)

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
var RawStdEncoding = StdEncoding.WithPadding(NoPadding)
// url base64编解码器 不用填充
var RawUrlEncoding = UrlEncoding.WithPadding(NoPadding)

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

// n个字节经过base64编码后的字节数组的长度(字节个数)
func (enc *Encoding) EncodedLen(n int) int {
	if enc.padChar == NoPadding {
		return (n * 8 + 5) / 6
	}
	return (n + 2) / 3 * 4
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
	// 此时8n/6不整除, 根据base64编码规则至少可以编码2个6个bit的字母
	val := uint(src[si+0]) << 16
	// 说明剩余2个字节
	if remain == 2 {
		val |= uint(src[si+1]) << 8
	}
	dst[di+0] = enc.encode[val>>18&0x3F]
	dst[di+1] = enc.encode[val>>12&0x3F]

	switch remain {
	case 1:
		// 需要进行尾部填充字符
		if enc.padChar != NoPadding {
			dst[di+2] = byte(enc.padChar)
			dst[di+3] = byte(enc.padChar)
		}
	case 2:
		dst[di+2] = enc.encode[val>>6&0x3F]
		if enc.padChar != NoPadding {
			dst[di+3] = byte(enc.padChar)
		}
	}
}

// 返回base编码后的字节数组对应的字符串
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}

// 返回base编码的字节数组(n个字节)经过解码后的长度(字节个数)
func (enc *Encoding) DecodedLen(n int) int {
	if enc.padChar == NoPadding {
		return n * 6 / 8
	}
	return n / 4 * 3
}

// 将base64编码的字节数组src解码到字节数组dst中
// 返回n表示解码后的字节个数
func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
	if len(src) == 0 {
		return 0, nil
	}
	si := 0
	for strconv.IntSize >= 64 && len(src) - si >= 8 && len(dst) - n >= 8 {
		if dn, ok := enc.decode64(src[si:]); ok {
			binary.BigEndian.PutUint64(dst[n:], dn)
			n += 6
			si += 8
		} else {
			var ninc int
			si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
			n += ninc
			if err != nil {
				return n, err
			}
		}
	}
	for len(src) - si >= 4 && len(dst) - n >= 4 {
		if dn, ok := enc.decode32(src[si:]); ok {
			binary.BigEndian.PutUint32(dst[n:], dn)
			n += 3
			si += 4
		} else {
			var ninc int
			si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
			n += ninc
			if err != nil {
				return n, err
			}
		}
	}
	for si < len(src) {
		var ninc int
		si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
		n += ninc
		if err != nil {
			return n, err
		}
	}
	return n, err
}

// 返回给定的base64编码的字符串解码后的字节数组
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

// 取base64编码后的字节数组的8个字节进行解码
func (enc *Encoding) decode64(src []byte) (dn uint64, ok bool) {
	var n uint64
	// 保证字节切片至少有8个字节
	// 如果没有直接panic出错
	_ = src[7]
	if n = uint64(enc.decodeMap[src[0]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 58
	if n = uint64(enc.decodeMap[src[1]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 52
	if n = uint64(enc.decodeMap[src[2]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 46
	if n = uint64(enc.decodeMap[src[3]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 40
	if n = uint64(enc.decodeMap[src[4]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 34
	if n = uint64(enc.decodeMap[src[5]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 28
	if n = uint64(enc.decodeMap[src[6]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 22
	if n = uint64(enc.decodeMap[src[7]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 16
	return dn, true
}

// 取base64编码后的字节数组的4个字节来进行解码
func (enc *Encoding) decode32(src []byte) (dn uint32, ok bool) {
	var n uint32
	_ = src[3]
	if n = uint32(enc.decodeMap[src[0]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 26
	if n = uint32(enc.decodeMap[src[1]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 20
	if n = uint32(enc.decodeMap[src[2]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 14
	if n = uint32(enc.decodeMap[src[3]]); n == 0xFF {
		return 0, false
	}
	dn |= n << 8
	return dn, true
}

func (enc *Encoding) decodeQuantum(dst, src []byte, si int) (nsi, n int, err error) {

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















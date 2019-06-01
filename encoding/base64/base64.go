package base64

// base64的编解码器
type Encoding struct {
	// 字母表(映射表)
	encode [64]byte
	// "逆"映射表
	decodeMap [256]byte
	// 尾部填充字符
	padChar rune
	strict bool
}

const (
	// 标准填充字符
	StdPadding rune = '='
	// 表示无填充字符
	NoPadding  rune = -1
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"


























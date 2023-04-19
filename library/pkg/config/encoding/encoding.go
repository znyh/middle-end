package encoding

import (
	"encoding/base64"
)

const (
	//BASE64字符表,不要有重复
	base64Table = "<>:;',./?~!@#$CDVWX%^&*ABYZabcghijklmnopqrstuvwxyz01EFGHIJKLMNOP"
)

var coder = base64.NewEncoding(base64Table)

/**
 * base64解密
 */
func Base64Decode(str string) (string, error) {
	by, err := coder.DecodeString(str)
	return string(by), err
}

func Base64Encode(str string) string {
	src := []byte(str)
	by := coder.EncodeToString(src)
	return by
}

package untils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const fixedSalt = "bitlongwallet7238baee9c2638664"

// generateSalt 生成一个指定长度的随机盐值
func generateMD5WithSalt(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input + fixedSalt))
	return hex.EncodeToString(hasher.Sum(nil))
}

// verifyChecksumWithSalt 验证字符串的校验码是否正确
func verifyChecksumWithSalt(originalString, checksum string) bool {
	expectedChecksum := generateMD5WithSalt(originalString)
	return checksum == expectedChecksum
}

func splitStringAndVerifyChecksum(extstring string) bool {
	parts := strings.Split(extstring, "_e_")
	if len(parts) != 2 {
		return false
	}
	originalString := parts[0]
	checksum := parts[1]
	return verifyChecksumWithSalt(originalString, checksum)

}

// 生成带校验值的扩展码
func GenerateExtMD5WithSalt(originalString string) string {
	expectedChecksum := generateMD5WithSalt(originalString)
	return originalString + "_e_" + expectedChecksum
}

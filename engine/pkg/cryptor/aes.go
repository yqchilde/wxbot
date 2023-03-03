package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// EncryptFilename 加密文件名
func EncryptFilename(key []byte, filename string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建随机的计数器
	counter := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, counter); err != nil {
		return "", err
	}

	// 使用CTR模式加密文件名
	stream := cipher.NewCTR(block, counter)
	ciphertext := make([]byte, len(filename))
	stream.XORKeyStream(ciphertext, []byte(filename))

	// 将计数器和密文合并为一个字节数组
	encrypted := make([]byte, aes.BlockSize+len(ciphertext))
	copy(encrypted[:aes.BlockSize], counter)
	copy(encrypted[aes.BlockSize:], ciphertext)

	// 对加密后的字节数组进行base64编码
	encoded := base64.RawURLEncoding.EncodeToString(encrypted)

	return encoded, nil
}

// DecryptFilename 解密文件名
func DecryptFilename(key []byte, encoded string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 对加密后的字符串进行base64解码
	encrypted, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	// 从解密后的字节数组中提取计数器和密文
	counter := encrypted[:aes.BlockSize]
	ciphertext := encrypted[aes.BlockSize:]

	// 使用CTR模式解密文件名
	stream := cipher.NewCTR(block, counter)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

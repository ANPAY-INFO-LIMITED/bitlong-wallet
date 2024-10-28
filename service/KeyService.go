package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/nbd-wtf/go-nostr"
	"math/rand"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"log"
)

const (
	keyId = "keyInfoId"
)

type PkInfo struct {
	Pubkey  string `json:"pubkey"`
	NpubKey string `json:"npubKey"`
}

func GetPrivateKey() (string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Println("err:", err)
	}
	if retrievedKey != nil {
		//PrivateKeyHex := fmt.Sprintf("%064x", retrievedKey.PrivateKey)
		//fmt.Println("PrivateKey:", PrivateKeyHex)
		//return PrivateKeyHex, nil
		privateKey := retrievedKey.PrivateKey
		return privateKey, nil
	}
	return "", fmt.Errorf("no key found")
}

func GenerateKeys(mnemonic, passphrase string) (string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Println("err:", err)
	}
	if retrievedKey != nil {
		publicKeyHex := fmt.Sprintf("%064x", retrievedKey.PublicKey)
		fmt.Println("pub:", publicKeyHex)
		return publicKeyHex, nil
	}
	fmt.Println("mnemonic:", mnemonic)
	seed := bip39.NewSeed(mnemonic, passphrase)
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", err
	}
	childKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", err
	}
	privKey, _ := btcec.PrivKeyFromBytes(childKey.Key)

	privateKeyHex := fmt.Sprintf("%064x", privKey.Serialize())
	publicKeyHex, err := getPublicKey(privateKeyHex)
	if err != nil {
		return "", err
	}
	fmt.Println("pub1:", publicKeyHex)
	err1 := saveDb(privateKeyHex, publicKeyHex)
	if err1 != nil {
		fmt.Println("err:", err1)
		return "", fmt.Errorf("failed to save keys: %s", err1)
	} // Assuming the bucket name is "Keys"
	return publicKeyHex, nil
}

func saveDb(privateKeyHex string, publicKeyHex string) error {
	db, err := InitDB() // 确保已经写了这个函数来初始化数据库和表
	if err != nil {
		fmt.Printf("Failed to initialize the database: %s\n", err)
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	keyStore := KeyStore{DB: db}

	keyInfo := KeyInfo{
		ID:         keyId, // This should be dynamically set or managed
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
	}
	key, err := keyStore.ReadKey(keyId)
	if err != nil {
		fmt.Printf("Failed to read data: %s\n", err)
	}
	if key == nil {
		err := keyStore.CreateOrUpdateKey(&keyInfo)
		if err != nil {
			fmt.Printf("Failed to save data: %s\n", err)
			return err
		}
	}
	return nil
}
func sign(privateKeyHex string, message string) (string, error) {
	var evt nostr.Event
	err := json.Unmarshal([]byte(message), &evt)
	if err != nil {
		return "", err
	}
	if err := evt.Sign(privateKeyHex); err != nil {
		return "", err
	}
	marshal, err := json.Marshal(evt)
	if err != nil {
		return "", err
	}
	return string(marshal), nil

}
func readDb() (*KeyInfo, error) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s\n", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	keyStore := &KeyStore{DB: db}
	// 调用 ReadKey 获取特定的密钥
	keyInfo, err := keyStore.ReadKey(keyId)
	if err != nil {
		log.Printf("Failed to read key %s: %s", keyId, err)
		return nil, err
	} else {
		//fmt.Printf("Key: %+v\n", keyInfo)
		return keyInfo, nil
	}
}

func SignMessage(message string) (string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
	}
	signInfo, err := sign(retrievedKey.PrivateKey, message)
	if err != nil {
		return "", err
	}

	fmt.Println(signInfo)
	return signInfo, nil
}
func getPublicKey(sk string) (string, error) {
	b, err := hex.DecodeString(sk)
	if err != nil {
		return "", err
	}
	_, pk := btcec.PrivKeyFromBytes(b)
	return hex.EncodeToString(schnorr.SerializePubKey(pk)), nil
}
func GetPublicRawKey() (string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
		return "", err
	}
	return retrievedKey.PublicKey, nil
}

func getNoStrAddress(pk string) (string, error) {
	compressedPubKeyBytes, err := hex.DecodeString(pk)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return "", err
	}
	// 将公钥编码为 Base58
	base58EncodedPubKey := base58.Encode(compressedPubKeyBytes)
	// 添加 nostr 协议所需的前缀
	nostrPubKey := "npub" + base58EncodedPubKey
	fmt.Println("Nostr address:", nostrPubKey)
	return nostrPubKey, nil
}

func GetPublicKey() (string, string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
		return "", "", err
	}
	publicKeyHex := fmt.Sprintf("%064X", retrievedKey.PublicKey)
	fmt.Println("publicKeyHex", publicKeyHex)
	address, err := getNoStrAddress(publicKeyHex)
	if err != nil {
		return "", "", err
	}
	return publicKeyHex, address, nil
}

// GetPublicKey 增强版获取公钥函数
func GetNewPublicKey() (string, string, error) {
	// 1. 从数据库读取密钥
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
		return "", "", err
	}

	// 2. 转换为64位十六进制格式
	publicKeyHex := fmt.Sprintf("%064X", retrievedKey.PublicKey)
	fmt.Println("publicKeyHex", publicKeyHex)

	// 3. 获取nostr地址
	nostrAddress, err := getNoStrAddress(publicKeyHex)
	if err != nil {
		return "", "", err
	}

	// 4. 加密nostr地址
	encryptedAddress, err := encryptNostrAddress(nostrAddress)
	if err != nil {
		return "", "", fmt.Errorf("encrypt address error: %v", err)
	}

	return publicKeyHex, encryptedAddress, nil
}

// insertRandomValues 在固定位置插入随机值
func insertRandomValues(publicKeyHex string) (string, error) {
	// 生成8位随机值
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomValue := hex.EncodeToString(randomBytes) // 8位十六进制

	// 每16个字符插入随机值
	var result strings.Builder
	for i := 0; i < len(publicKeyHex); i += 16 {
		end := i + 16
		if end > len(publicKeyHex) {
			end = len(publicKeyHex)
		}
		result.WriteString(publicKeyHex[i:end])
		if end < len(publicKeyHex) {
			result.WriteString("_")
			result.WriteString(randomValue)
		}
	}

	return result.String(), nil
}

// encryptPublicKey AES加密公钥
func encryptNostrAddress(address string) (string, error) {
	// 1. 插入随机值
	addressWithRandom, err := insertRandomValues(address)
	if err != nil {
		return "", err
	}
	// 2. AES加密
	return aesEncrypt(addressWithRandom)
}
func aesEncrypt(data string) (string, error) {
	// 1. 验证输入
	if len(data) == 0 {
		return "", fmt.Errorf("empty data")
	}

	// 2. AES密钥 (32字节 for AES-256)
	key := []byte("YourAESKey32BytesLongForSecurity")
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: must be 32 bytes")
	}

	// 3. 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher error: %v", err)
	}

	// 4. 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("generate IV error: %v", err)
	}

	// 5. PKCS7填充（使用正确的块大小）
	paddedData := pkcs7Pad([]byte(data), aes.BlockSize)

	// 6. 创建加密文本缓冲区
	ciphertext := make([]byte, len(paddedData))

	// 7. 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	// 8. 组合IV和密文
	combined := make([]byte, len(iv)+len(ciphertext))
	copy(combined, iv)
	copy(combined[len(iv):], ciphertext)

	// 9. 返回十六进制编码的结果
	return hex.EncodeToString(combined), nil
}

// PKCS7填充函数
func pkcs7Pad(data []byte, blockSize int) []byte {
	// 1. 验证块大小
	if blockSize <= 0 || blockSize > 256 {
		panic("invalid block size")
	}

	// 2. 计算需要填充的长度
	padding := blockSize - len(data)%blockSize

	// 3. 创建填充数据
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	// 4. 添加填充
	return append(data, padText...)
}
func GetJsonPublicKey() (string, error) {
	var pkInfo PkInfo
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
		return "", err
	}
	publicKeyHex := fmt.Sprintf("%064X", retrievedKey.PublicKey)
	fmt.Println("publicKeyHex", publicKeyHex)
	address, err := getNoStrAddress(publicKeyHex)
	if err != nil {
		return "", err
	}
	pkInfo.Pubkey = publicKeyHex
	pkInfo.NpubKey = address
	marshal, err := json.Marshal(pkInfo)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

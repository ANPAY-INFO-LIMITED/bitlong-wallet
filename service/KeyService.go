package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/nbd-wtf/go-nostr"
	"golang.org/x/crypto/pbkdf2"
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

func readDb_Fixed() (*KeyInfo, error) {
	db, err := InitDB1()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s\n", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	keyStore := &KeyStore{DB: db}
	keyInfo, err := keyStore.ReadKey(keyId)
	if err != nil {
		log.Printf("Failed to read key %s: %s", keyId, err)
		return nil, err
	} else {
		return keyInfo, nil
	}
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
	keyInfo, err := keyStore.ReadKey(keyId)
	if err != nil {
		log.Printf("Failed to read key %s: %s", keyId, err)
		return nil, err
	} else {
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
	base58EncodedPubKey := base58.Encode(compressedPubKeyBytes)
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

func GetNewPublicKey() (string, string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
		return "", "", err
	}

	publicKeyHex := fmt.Sprintf("%064X", retrievedKey.PublicKey)
	fmt.Println("publicKeyHex", publicKeyHex)

	nostrAddress, err := getNoStrAddress(publicKeyHex)
	if err != nil {
		return "", "", err
	}

	encryptedAddress, err := encryptNostrAddress(nostrAddress)
	if err != nil {
		return "", "", fmt.Errorf("encrypt address error: %v", err)
	}

	return publicKeyHex, encryptedAddress, nil
}

func insertRandomValues(publicKeyHex string) (string, error) {
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomValue := hex.EncodeToString(randomBytes) // 8位十六进制

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

func encryptNostrAddress(address string) (string, error) {
	addressWithRandom, err := insertRandomValues(address)
	if err != nil {
		return "", err
	}
	return aesEncrypt(addressWithRandom)
}
func aesEncrypt(data string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty data")
	}

	key := []byte("YourAESKey32BytesLongForSecurity")
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher error: %v", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("generate IV error: %v", err)
	}

	paddedData := pkcs7Pad([]byte(data), aes.BlockSize)

	ciphertext := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	combined := make([]byte, len(iv)+len(ciphertext))
	copy(combined, iv)
	copy(combined[len(iv):], ciphertext)

	return hex.EncodeToString(combined), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	if blockSize <= 0 || blockSize > 256 {
		panic("invalid block size")
	}

	padding := blockSize - len(data)%blockSize

	padText := bytes.Repeat([]byte{byte(padding)}, padding)

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
func decryptNew(cipherText string, key []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %w", err)
	}

	if len(decoded) < aes.BlockSize {
		return "", errors.New("解码后的数据长度不足16字节")
	}

	iv := decoded[:aes.BlockSize]
	encrypted := decoded[aes.BlockSize:]

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("无效的密钥长度")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建AES密钥失败: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	if len(decrypted) == 0 {
		return "", errors.New("解密后的数据为空")
	}
	padding := int(decrypted[len(decrypted)-1])

	if padding < 1 || padding > aes.BlockSize || padding > len(decrypted) {
		return "", errors.New("无效的填充长度")
	}

	for i := 0; i < padding; i++ {
		if decrypted[len(decrypted)-1-i] != byte(padding) {
			return "", errors.New("填充验证失败")
		}
	}

	return string(decrypted[:len(decrypted)-padding]), nil
}

func BuildDecrypt(encryptedDeviceID, saltBase64 string) (string, error) {
	password := []byte("thisisaverysecretkey1234567890")
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		fmt.Println("解码盐值失败:", err)
		return "", err
	}
	key := pbkdf2.Key(password, salt, 10000, 32, sha256.New)
	decryptedID, err := decryptNew(encryptedDeviceID, key)
	if err != nil {
		fmt.Println("解密失败:", err)
		return "", err
	}
	return decryptedID, nil
}
func GetExistPublicKey() (string, string, error) {
	retrievedKey, err := readDb()
	if err != nil {
		fmt.Printf("err is :%v\n", err)
		return "", "", err
	}

	publicKeyHex := fmt.Sprintf("%064X", retrievedKey.PublicKey)
	fmt.Println("publicKeyHex", publicKeyHex)

	nostrAddress, err := getNoStrAddress(publicKeyHex)
	if err != nil {
		return "", "", err
	}

	encryptedAddress, err := encryptNostrAddress(nostrAddress)
	if err != nil {
		return "", "", fmt.Errorf("encrypt address error: %v", err)
	}

	return publicKeyHex, encryptedAddress, nil
}

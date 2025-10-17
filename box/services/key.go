package services

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"github.com/wallet/box/db"
	"github.com/wallet/box/models"
	"gorm.io/gorm"
	"os"
)

const (
	nullSeedNPub = "npubzGfdXpacbK5B9FEX1M3Ud29vt1owYvzLAtN5vgTf7xS1vZunBZgZf555o8HUbtUz73uZMvHdPwcsk4h3X4f9xx7"
)

func UpdateKey() error {
	seed, err := getSeed()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logrus.Infoln("seed not found, skipping key update")
			return nil
		}

		return errors.Wrap(err, "getSeed")
	}
	if seed == "" {
		logrus.Infoln("seed is empty, skipping key update")
		return nil
	}

	priKey, pubKey, nPub, err := getKeys(seed)
	if err != nil {
		return errors.Wrap(err, "getKeys")
	}
	if nPub == nullSeedNPub {
		return errors.Wrap(invalidNPub, "nPub is nullSeedNPub")
	}

	tx := db.Sqlite().Begin()

	var k models.Key
	err = tx.Model(&models.Key{}).First(&k).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = tx.Model(&models.Key{}).
				Create(&models.Key{
					PriKey: priKey,
					PubKey: pubKey,
					NPub:   nPub,
				}).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "tx.Model(&models.Key{}).Create")
			}
		} else {
			tx.Rollback()
			return errors.Wrap(err, "tx.Model(&models.Key{}).First")
		}
	} else {
		logrus.Infof("priKey: %s, pubKey: %s, nPub: %s", priKey, pubKey, nPub)

		if priKey != k.PriKey || pubKey != k.PubKey || nPub != k.NPub {
			if k.NPub == nullSeedNPub {
				priKey = ""
				pubKey = ""
				nPub = ""
			}

			err = tx.Model(&models.Key{}).
				Where("id = ?", k.ID).
				Updates(map[string]any{
					"pri_key": priKey,
					"pub_key": pubKey,
					"n_pub":   nPub,
				}).Error
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "tx.Model(&models.Key{}).Save")
			}
		} else {
			return tx.Rollback().Error
		}
	}

	return tx.Commit().Error
}

func GetNPub() (string, error) {
	var k models.Key
	tx := db.Sqlite().Begin()
	err := tx.Model(&models.Key{}).First(&k).Error
	if err != nil {
		tx.Rollback()
		return "", errors.Wrap(err, "tx.Model(&models.Key{}).First")
	}
	tx.Rollback()
	return k.NPub, nil
}

func getKeys(mnemonic string) (string, string, string, error) {
	priKey, pubKey, err := genKeys(mnemonic)
	if err != nil {
		return "", "", "", errors.Wrap(err, "genKeys")
	}
	nPub, err := nPubKey(pubKey)
	if err != nil {
		return "", "", "", errors.Wrap(err, "nPubKey")
	}
	return priKey, pubKey, nPub, nil
}

func genKeys(mnemonic string) (string, string, error) {
	passphrase := ""
	seed := bip39.NewSeed(mnemonic, passphrase)
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", "", errors.Wrap(err, "bip32.NewMasterKey")
	}
	childKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", "", errors.Wrap(err, "masterKey.NewChildKey")
	}
	priKey, _ := btcec.PrivKeyFromBytes(childKey.Key)

	privateKeyHex := fmt.Sprintf("%064x", priKey.Serialize())

	publicKeyHex, err := getPublicKey(privateKeyHex)
	if err != nil {
		return "", "", errors.Wrap(err, "getPublicKey")
	}

	return privateKeyHex, publicKeyHex, nil
}

func getPublicKey(sk string) (string, error) {
	b, err := hex.DecodeString(sk)
	if err != nil {
		return "", errors.Wrap(err, "hex.DecodeString")
	}
	_, pk := btcec.PrivKeyFromBytes(b)
	return hex.EncodeToString(schnorr.SerializePubKey(pk)), nil
}

func getNoStrAddr(pk string) (string, error) {
	compressedPubKeyBytes, err := hex.DecodeString(pk)
	if err != nil {
		return "", errors.Wrap(err, "hex.DecodeString")
	}
	base58EncodedPubKey := base58.Encode(compressedPubKeyBytes)
	noStrPubKey := "npub" + base58EncodedPubKey
	return noStrPubKey, nil
}

func nPubKey(pubKey string) (string, error) {
	publicKeyHex := fmt.Sprintf("%064X", pubKey)
	addr, err := getNoStrAddr(publicKeyHex)
	if err != nil {
		return "", errors.Wrap(err, "getNoStrAddr")
	}
	return addr, nil
}

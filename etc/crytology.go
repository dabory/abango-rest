// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
)

func RsaPrPbPair(keyLeng int) ([]byte, []byte) {

	prKey, _ := rsa.GenerateKey(rand.Reader, keyLeng)
	pbKey := &prKey.PublicKey

	prBytes := x509.MarshalPKCS1PrivateKey(prKey)
	prMem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: prBytes,
	})

	// Bytes: 에 직접 넣으면 런타임에서 에러남.(중요!!)
	pbBytes, err := x509.MarshalPKIXPublicKey(pbKey)
	if err != nil {
		fmt.Println(err)
	}
	pbMem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pbBytes,
	})

	return prMem, pbMem
}

func RsaSignature(prKey []byte, msg []byte) ([]byte, error) { // msg 245=(256-11)bytes 이하

	block, _ := pem.Decode(prKey)
	if block == nil {
		// fmt.Println("Error: pem.Decode in MySignature")
		return nil, MyErr("QROPBDHCF-pem.Decode", nil, false)
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, MyErr("QERDVAEQYG-x509.ParsePKCS1PrivateKey", err, false)
	}

	sign, err := rsa.SignPKCS1v15(nil, priv, crypto.Hash(0), msg)
	if err != nil {
		return nil, MyErr("CZGSRFVA-rsa.SignPKCS1v15", err, false)
	}
	return sign, nil

}

func RsaOriginal(pubKey []byte, msg []byte) ([]byte, error) {

	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, MyErr("ADYERFBJ-pem.Decode", nil, false)
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, MyErr("MNCHSTRGB-x509.ParsePKIXPublicKey", err, false)
	}
	pbKey := pubInterface.(*rsa.PublicKey)

	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(msg)
	e := big.NewInt(int64(pbKey.E))
	c.Exp(m, e, pbKey.N)
	out := c.Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}
	return out[skip:], nil
}

func RsaPbEncrypt(publicKey []byte, msg []byte) ([]byte, error) {
	origData := msg
	block, _ := pem.Decode(publicKey)

	if block == nil {
		return nil, MyErr("QEROBVSRAE-pem.Decode", nil, false)
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, MyErr("KDGBSAERT-x509.ParsePKIXPublicKey", err, false)
	}
	pub := pubInterface.(*rsa.PublicKey)

	label := []byte("")
	sha256hash := sha256.New()
	enBytes, err := rsa.EncryptOAEP(sha256hash, rand.Reader, pub, origData, label)
	if err != nil {
		return nil, MyErr("rsa.EncryptOAEP", err, false)
	}

	return enBytes, nil
}

func RsaPrDecrypt(privateKey []byte, msg []byte) ([]byte, error) {

	ciphertext := msg
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, MyErr("USFBSHWER-pem.Decode", nil, false)
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, MyErr("MTVZQFGAEW-x509.ParsePKCS1PrivateKey", err, false)
	}

	label := []byte("")
	sha256hash := sha256.New()
	deBytes, err := rsa.DecryptOAEP(sha256hash, rand.Reader, priv, ciphertext, label)
	if err != nil {
		return nil, MyErr("QGOTRVBAE-rsa.DecryptOAEP", err, false)
	}
	return deBytes, nil
}

func AddBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func RemoveBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func AesPad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func AesUnpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("OYNSTBZE-unpad error. This could happen when incorrect myEncryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func AesEncrypt(key []byte, text []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("WERFAV-Error: NewCipher in myEncrypt - " + err.Error())
		return nil, err
	}

	msg := AesPad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, MyErr("-NYTZVCEF-io.ReadFull", err, false)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)
	finalMsg := RemoveBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))

	return []byte(finalMsg), nil
}

func AesDecrypt(key []byte, text []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, MyErr("YGBARZDF-aes.NewCipher", err, false)
	}
	decodedMsg, err := base64.URLEncoding.DecodeString(AddBase64Padding(string(text)))
	if err != nil {
		return nil, MyErr("VSGRERGBEW-base64.URLEncoding.DecodeString, Possibley Decryption string is too long", err, false)
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return nil, MyErr("MHETYGVBA-aes.BlockSize-blocksize must be multipe of decoded message length", err, false)
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := AesUnpad(msg)
	if err != nil {
		return nil, MyErr("WERQFAERGQ-AesUnpad", err, false)
	}
	return unpadMsg, nil
}

func Sha256Hash(data []byte, leng int) []byte {
	hash := sha256.New() //SHA-3 규격임.
	hash.Write(data)

	mdStr := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	rtn := ""
	if leng == 0 {
		rtn = mdStr
	} else {
		rtn = mdStr[10 : 10+leng]
	}
	return []byte(rtn)
}

func Aes256Encrypt(key []byte, nonce []byte, plaintext []byte) ([]byte, error) {
	// The key argument should be the AES key, either 16 or 32 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, MyErr("YBBAERFRYY-NewCipher", err, false)
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// nonce := make([]byte, 12) // Do not change 12
	// if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	// 	return nil,  myErr("io.ReadFull", err)
	// }

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, MyErr("ERGVSER-cipher.NewGCM", err, false)
	}

	text := aesgcm.Seal(nil, nonce, plaintext, nil)
	return text, nil
}

func Aes256Decrypt(key []byte, nonce []byte, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, MyErr("NSTRFGBSAF-NewCipher", err, false)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, MyErr("PQWKKLVASD-cipher.NewGCM", err, false)
	}

	plaintext, err := aesgcm.Open(nil, nonce, text, nil)
	if err != nil {
		return nil, MyErr("POQWEIRUNVAIK-aesgcm.Open", err, false)
	}
	return plaintext, nil
}

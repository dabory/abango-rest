package locals

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"

	"weberp-go/locals/gosodium/cryptobox"

	e "github.com/dabory/abango-rest/etc"
)

func PkeyDecrypt(encr64 string, keyPair64 string) ([]byte, error) {

	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
	if err != nil {
		return nil, e.LogErr("ertjhdssfwew", "Decryption Failure-1 ", err)
	}

	sKey, pKey, err := cryptobox.CryptoBoxGetSecretPublicKeyFrom(keyPair)
	if err != nil {
		return nil, e.LogErr("034uldhjalrse", "Decryption Failure-2 ", err)
	}

	decodedStr, err := base64.StdEncoding.DecodeString(encr64)
	if err != nil {
		return nil, e.LogErr("32rfwr32rfw3", "Decryption Failure-3 ", err)
	}

	decryptedBytes, boxRet := cryptobox.CryptoBoxSealOpen(decodedStr, pKey, sKey)
	if boxRet != 0 {
		return nil, e.LogErr("mcnbxkajhr3eih", "Decryption Failure-4, boxRet:"+strconv.Itoa(boxRet), nil)
	}
	return decryptedBytes, nil
}

func PkeyEncrypt(msg string, keyPair64 string) (string, error) {

	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
	if err != nil {
		return "", e.LogErr("2903uljslw", "Encrypt Failure-1 ", err)
	}

	_, pKey, err := cryptobox.CryptoBoxGetSecretPublicKeyFrom(keyPair)
	if err != nil {
		return "", e.LogErr("mcnoiaoruec", "Encrypt Failure-2 ", err)
	}

	EncryptedBytes, boxRet := cryptobox.CryptoBoxSeal([]byte(msg), pKey)
	if boxRet != 0 {
		return "", e.LogErr("8u3h82f0ud", "Encrypt Failure-4, boxRet:"+strconv.Itoa(boxRet), nil)
	}
	EncryptedBase64 := base64.StdEncoding.EncodeToString(EncryptedBytes)
	return EncryptedBase64, nil
}

func DbrPasswd(password string, salt string) string {
	salt16 := DbrSaltBase(salt, 16)
	var passwordBytes = []byte(password)
	var sha256Hasher = sha256.New()

	passwordBytes = append(passwordBytes, salt16...)
	sha256Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha256Hasher.Sum(nil)
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

func DbrHashedIndex(target string) string {
	//!!중요: salt는 16char에서만 작동된다. hash 값은 44 char나오지만 32char로 잘라서 쓴다.
	// fmt.Println("hash_full_length:", DbrPasswd(target, "$$hashed_index$$"))
	return DbrPasswd(target, "$$hashed_index$$")[0:32]
}

func DbrCompare(hashedPassword, currPassword string, salt string) bool {
	// fmt.Println("salt:", salt)
	// fmt.Println("currPassword:", currPassword)
	var currPasswordHash = DbrPasswd(currPassword, salt)
	// fmt.Println("currPasswordHash:", currPasswordHash)
	// fmt.Println("hashedPassword:", hashedPassword)
	return hashedPassword == currPasswordHash
}

func DbrSaltBase(salt string, saltSize int) []byte { //어떤 사이즈라도 16byte의 Base64로 변경
	tmp := []byte(salt)
	salt64 := base64.StdEncoding.EncodeToString(tmp)
	return []byte(salt64[4 : saltSize+4])
}

// if keysize is 16bytes * 8bits -> 128
// if keysize is 32bytes * 8bits -> 256
// Encrypt-Decript는 plaintext가 16bytes 밖에는 지원하지 않는다 따라서
// MyAesEncrypt를 사용한다.
func MyAesEncrypt(key []byte, text []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New(e.FuncRunErr("odvjkwei3", e.CurrFuncName()+" "+err.Error()))
	}

	msg := Pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.New(e.FuncRunErr("ls0ue3so", e.CurrFuncName()+" "+err.Error()))
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))

	return []byte(finalMsg), nil
}

func MyAesDecrypt(key []byte, text []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New(e.FuncRunErr("3do8awe", e.CurrFuncName()+" "+err.Error()))
	}
	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(string(text)))
	if err != nil {
		return nil, errors.New(e.FuncRunErr("mkshewjd", e.CurrFuncName()+" "+err.Error()))
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return nil, errors.New(e.FuncRunErr("mskoeuwid", e.CurrFuncName()+" "+err.Error()))
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return nil, errors.New(e.FuncRunErr("012bsoo832d", e.CurrFuncName()+" "+err.Error()))
	}
	return unpadMsg, nil
}

func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New(e.FuncRunErr("unpad error. This could happen when incorrect MyAesEncryption key is used", e.CurrFuncName()))
	}
	return src[:(length - unpadding)], nil
}

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}
	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

package etc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/dabory/abango-rest/gosodium/cryptobox"
)

func KeyPairGenerate() (string, error) {

	publicKey, secretKey, boxRet := cryptobox.CryptoBoxKeyPair()
	if boxRet != 0 {
		return "", LogErr("ertjhdswew", "CryptoBoxGenerateKeyPair", errors.New("dljdf"))
	}

	kpBytes := make([]byte, 64)
	copy(kpBytes[:32], publicKey) // Extracting the secret key part
	copy(kpBytes[32:], secretKey) // Appending the public key part
	return base64.StdEncoding.EncodeToString(kpBytes), nil
}

func PkeyDecrypt(encr64 string, keyPair64 string) ([]byte, error) {

	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
	if err != nil {
		return nil, LogErr("ertjhdssw", FuncNameErr()+"Failure-1 ", err)
	}

	sKey, pKey, err := cryptobox.CryptoBoxGetSecretPublicKeyFrom(keyPair)
	if err != nil {
		return nil, LogErr("034hjalrse", FuncNameErr()+"Failure-2 ", err)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encr64)
	if err != nil {
		return nil, LogErr("32rfww3", FuncNameErr()+"Failure-3 ", err)
	}

	decryptedBytes, boxRet := cryptobox.CryptoBoxSealOpen(decodedBytes, pKey, sKey)
	if boxRet != 0 {
		return nil, LogErr("mcnbxkajhr3eih", FuncNameErr()+"boxRet:"+strconv.Itoa(boxRet), nil)
	}
	return decryptedBytes, nil
}

func PkeyEncrypt(msg string, keyPair64 string) (string, error) {

	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
	if err != nil {
		return "", LogErr("2903uljslw", FuncNameErr()+"Failure-1 ", err)
	}

	_, pKey, err := cryptobox.CryptoBoxGetSecretPublicKeyFrom(keyPair)
	if err != nil {
		return "", LogErr("mcnoiaoc", FuncNameErr()+"Failure-2 ", err)
	}

	EncryptedBytes, boxRet := cryptobox.CryptoBoxSeal([]byte(msg), pKey)
	if boxRet != 0 {
		return "", LogErr("8u3h82f0", FuncNameErr()+"boxRet:"+strconv.Itoa(boxRet), nil)
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
	return base64.URLEncoding.EncodeToString(hashedPasswordBytes)
}

// md5는 간단한 Device Hash 같은 간단한 hash 이용하기 위해서
func Md5Hashed(target string, size int) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(target)))[0:size]
}

func DbrHashed(target string, size int) string {
	//!!중요: salt는 16char에서만 작동된다. hash 값은 44 char나오지만 size로 잘라서 쓴다.
	return DbrPasswd(target, "$$email_hashed$$")[0:size]
}

func DbrCompare(hashedPassword, currPassword string, salt string) bool {
	var currPasswordHash = DbrPasswd(currPassword, salt)
	return hashedPassword == currPasswordHash
}

func DbrSaltBase(salt string, saltSize int) []byte { //어떤 사이즈라도 16byte의 Base64로 변경
	tmp := []byte(salt)
	salt64 := base64.StdEncoding.EncodeToString(tmp)
	return []byte(salt64[4 : saltSize+4])
}

// if keysize is 16bytes * 8bits -> 128
// if keysize is 32bytes * 8bits -> 256
// Encrypt-Decrypt는 plaintext가 16bytes 밖에는 지원하지 않는다 따라서 MyAesEncrypt를 사용한다.
func MyAesEncrypt(key []byte, text []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New(FuncRunErr("odvjkwei3", FuncNameErr()+" "+err.Error()))
	}

	msg := Pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.New(FuncRunErr("ls0ue3so", FuncNameErr()+" "+err.Error()))
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))

	return []byte(finalMsg), nil
}

func MyAesDecrypt(key []byte, text []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New(FuncRunErr("3do8awe", FuncNameErr()+" "+err.Error()))
	}
	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(string(text)))
	if err != nil {
		return nil, errors.New(FuncRunErr("mkshewjd", FuncNameErr()+" "+err.Error()))
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return nil, errors.New(FuncRunErr("mskoeuwid", FuncNameErr()+" "+err.Error()))
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return nil, errors.New(FuncRunErr("012bsoo832d", FuncNameErr()+" "+err.Error()))
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
		return nil, errors.New(FuncRunErr("unpad error. This could happen when incorrect MyAesEncryption key is used", FuncNameErr()))
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

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

const ( //Secutiry
	SaltyKeyPairPrefix string = "@_@_"
	BelovedPass        string = "20150721-20200102" //Do NOT Change
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

	var err error
	if keyPair64, err = DecryptKeyPair(keyPair64); err != nil {
		return nil, LogErr("ertjhdssw", FuncNameErr()+"Failure-0 ", err)
	}

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

func DecryptKeyPair(keyPair64 string) (string, error) {
	// 1) prefix 없는 경우 방어
	if !strings.HasPrefix(keyPair64, SaltyKeyPairPrefix) {
		return keyPair64, nil
	}

	// 2) prefix 정확히 한 번만 제거
	deSalted := strings.TrimPrefix(keyPair64, SaltyKeyPairPrefix)
	fmt.Println("deSalted: ", deSalted)
	fmt.Println("len-deSalted: ", len(deSalted))
	aesKey := DeriveAesKey(BelovedPass)
	decrKeyPair, err := AesGcmDecrypt(aesKey, deSalted)
	if err != nil {
		return "", LogErr("salty-dec", FuncNameErr()+"AesGcmDecrypt failed", err)
	}

	return string(decrKeyPair), nil
}

func EncryptKeyPair(keyPair64 string) (string, error) {
	aesKey := DeriveAesKey(BelovedPass)
	encrKeyPair, err := AesGcmEncrypt(aesKey, []byte(keyPair64))
	if err != nil {
		return "", LogErr("salty-enc", FuncNameErr()+"AesGcmEncrypt failed", err)
	}

	return SaltyKeyPairPrefix + encrKeyPair, nil
}

func PkeyEncrypt(msg string, keyPair64 string) (string, error) {

	var err error
	if keyPair64, err = DecryptKeyPair(keyPair64); err != nil {
		return "", LogErr("ertjhds6w", FuncNameErr()+"Failure-0 ", err)
	}

	fmt.Println("keyPair64:", keyPair64)

	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
	if err != nil {
		return "", LogErr("2903uljslw", FuncNameErr()+"Failure-1 ", err)
	}
	// fmt.Println("keyPair:", string(keyPair))

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

func DbrPasswd(password string, salt string) string { // 반값이면 빈값을 리컨한다.
	if password == "" {
		return ""
	}

	salt16 := DbrSaltBase(salt, 16)
	var passwordBytes = []byte(password)
	var sha256Hasher = sha256.New()

	passwordBytes = append(passwordBytes, salt16...)
	sha256Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha256Hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashedPasswordBytes)
}

// func DbrPasswd(password string, salt string) string {
// 	salt16 := DbrSaltBase(salt, 16)
// 	var passwordBytes = []byte(password)
// 	var sha256Hasher = sha256.New()

// 	passwordBytes = append(passwordBytes, salt16...)
// 	sha256Hasher.Write(passwordBytes)

// 	var hashedPasswordBytes = sha256Hasher.Sum(nil)
// 	return base64.URLEncoding.EncodeToString(hashedPasswordBytes)
// }

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

// AesGcmEncrypt : AES-GCM 기반 암호화
// key: 32byte (sha256 결과)
// plainText: 암호화할 원문
// 리턴: base64(nonce||ciphertext||tag)
// 받는건 최소 128 Char 이상이여야 하는 것 같다.
func AesGcmEncrypt(key []byte, plainText []byte) (string, error) {

	if len(plainText) == 0 { // 값이 없어도 에러를 내지 않는다.
		return "", nil
	}
	// 1. key와 plainText 빈값 및 유효성 체크
	if len(key) == 0 {
		return "", LogErr("key_empty", FuncNameErr()+": ", errors.New("key is empty"))
	}

	// AES 키 길이 검증 (16, 24, 32 바이트)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", LogErr("key_invalid", FuncNameErr()+": ", errors.New("invalid key size"))
	}

	// 2. Cipher 블록 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", LogErr("3do8awe", FuncNameErr()+":NewCipher ", err)
	}

	// 3. GCM 모드 인스턴스 생성
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", LogErr("gcm_init_err", FuncNameErr()+":NewGCM ", err)
	}

	// 4. Nonce 생성 (임의의 난수)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", LogErr("ls0ue3so", FuncNameErr()+":ReadFull ", err)
	}

	// 5. 암호화 실행 (Seal)
	// nonce를 첫 번째 인자로 전달하여 암호문 앞에 자동으로 붙게 함 (nonce + ciphertext + tag)
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	// 6. 결과를 Base64 문자열로 인코딩 후 []byte로 반환
	encBase64 := base64.StdEncoding.EncodeToString(cipherText)

	return encBase64, nil
}

// func AesGcmEncrypt(key []byte, plainText []byte) ([]byte, error) {

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, errors.New("AesGcmEncrypt NewCipher: " + err.Error())
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, errors.New("AesGcmEncrypt NewGCM: " + err.Error())
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, errors.New("AesGcmEncrypt nonce: " + err.Error())
// 	}

// 	// gcm.Seal: nonce + encrypted + tag (AEAD)
// 	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

// 	encBase64 := base64.StdEncoding.EncodeToString(cipherText)
// 	return []byte(encBase64), nil
// }

// AesGcmDecrypt : AES-GCM 기반 복호화
// encBase64: base64(nonce||ciphertext||tag)
func AesGcmDecrypt(key []byte, encBase64 string) (string, error) {

	if len(encBase64) == 0 { // 값이 없어도 에러를 내지 않는다.
		return "", nil
	}
	// 1. key와 encBase64 빈값 및 유효성 체크
	if len(key) == 0 {
		return "", LogErr("key_empty", FuncNameErr()+": ", errors.New("key is empty"))
	}

	// AES 키 길이 검증 (GCM에서도 AES 키 규격은 동일함)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", LogErr("key_invalid", FuncNameErr()+": ", errors.New("invalid key size"))
	}
	// 2. Base64 디코딩
	cipherBlob, err := base64.StdEncoding.DecodeString(encBase64)
	if err != nil {
		return "", LogErr("mkshewjd", FuncNameErr()+":base64 decode ", err)
	}

	// 3. Cipher 블록 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", LogErr("3do8awe", FuncNameErr()+":NewCipher ", err)
	}

	// 4. GCM 모드 인스턴스 생성
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", LogErr("gcm_init_err", FuncNameErr()+":NewGCM ", err)
	}

	// 5. 데이터 길이 검증 (Nonce 포함 여부)
	nonceSize := gcm.NonceSize()
	if len(cipherBlob) < nonceSize {
		return "", LogErr("mskoeuwid", FuncNameErr()+": ", errors.New("cipher too short"))
	}

	// 6. Nonce와 CipherText 분리
	nonce := cipherBlob[:nonceSize]
	cipherText := cipherBlob[nonceSize:]

	// 7. 복호화 및 검증(Open)
	byteMsg, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", LogErr("012bsoo832d", FuncNameErr()+":open ", err)
	}

	return string(byteMsg), nil
}

// func AesGcmDecrypt(key []byte, encBase64 []byte) ([]byte, error) {

// 	// base64 decode -> []byte blob
// 	cipherBlob, err := base64.StdEncoding.DecodeString(string(encBase64))
// 	if err != nil {
// 		return nil, errors.New("AesGcmDecrypt base64 decode: " + err.Error())
// 	}

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, errors.New("AesGcmDecrypt NewCipher: " + err.Error())
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, errors.New("AesGcmDecrypt NewGCM: " + err.Error())
// 	}

// 	nonceSize := gcm.NonceSize()
// 	if len(cipherBlob) < nonceSize {
// 		return nil, errors.New("AesGcmDecrypt cipher too short")
// 	}

// 	nonce := cipherBlob[:nonceSize]
// 	cipherText := cipherBlob[nonceSize:]

// 	byteMsg, err := gcm.Open(nil, nonce, cipherText, nil)
// 	if err != nil {
// 		return nil, errors.New("AesGcmDecrypt open: " + err.Error())
// 	}

// 	// plainText 는 []byte 그대로 반환
// 	return byteMsg, nil
// }

// DeriveAesKey : passphrase 로부터 32byte AES 키 생성
func DeriveAesKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase)) // 32bytes fixed
	return hash[:]                            // slice 로 변환
}

// 이것 여기에만 쓰는 것이니까 제거해 볼것.
// locals/gate-token-owner-key-related.go:94:20: undefined: e.MyAesDecrypt
func MyAesEncrypt(key []byte, text []byte) ([]byte, error) {
	// 1. key와 text 빈값 및 유효성 체크 추가
	if len(key) == 0 {
		return nil, LogErr("odvjkwei3", FuncNameErr()+": ", errors.New("key is empty"))
	}

	// AES 키 길이는 16(AES-128), 24(AES-192), 32(AES-256) 바이트여야 합니다.
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, LogErr("odvjkwei3", FuncNameErr()+": ", errors.New("invalid key size"))
	}

	if len(text) == 0 {
		return nil, LogErr("odvjkwei3", FuncNameErr()+": ", errors.New("text is empty"))
	}

	// 2. 기존 로직 시작
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, LogErr("odvjkwei3", FuncNameErr()+":NewCipher ", err)
	}

	msg := Pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, LogErr("ls0ue3so", FuncNameErr()+":ReadFull ", err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))

	return []byte(finalMsg), nil
}

func MyAesDecrypt(key []byte, text []byte) ([]byte, error) {
	// 1. key 빈값 및 유효성 체크
	if len(key) == 0 {
		return nil, LogErr("key_empty", FuncNameErr()+": ", errors.New("key is empty"))
	}

	// AES 키 길이 검증
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, LogErr("key_invalid", FuncNameErr()+": ", errors.New("invalid key size"))
	}

	// 2. text 빈값 체크
	if len(text) == 0 {
		return nil, LogErr("text_empty", FuncNameErr()+": ", errors.New("text is empty"))
	}

	// 3. Cipher 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, LogErr("3do8awe", FuncNameErr()+": ", err)
	}

	// 4. Base64 디코딩
	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(string(text)))
	if err != nil {
		return nil, LogErr("mkshewjd", FuncNameErr()+": ", err)
	}

	// 5. 데이터 길이 검증 (IV 포함 여부)
	if len(decodedMsg) < aes.BlockSize {
		return nil, LogErr("mskoeuwid", FuncNameErr()+": ", errors.New("ciphertext too short"))
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	// 6. 복호화 실행
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	// 7. 패딩 제거
	unpadMsg, err := Unpad(msg)
	if err != nil {
		return nil, LogErr("012bsoo832d", FuncNameErr()+": ", err)
	}

	return unpadMsg, nil
}

// func MyAesEncrypt(key []byte, text []byte) ([]byte, error) {

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, errors.New(FuncRunErr("odvjkwei3", FuncNameErr()+" "+err.Error()))
// 	}

// 	msg := Pad(text)
// 	ciphertext := make([]byte, aes.BlockSize+len(msg))
// 	iv := ciphertext[:aes.BlockSize]
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return nil, errors.New(FuncRunErr("ls0ue3so", FuncNameErr()+" "+err.Error()))
// 	}

// 	cfb := cipher.NewCFBEncrypter(block, iv)
// 	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)
// 	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))

// 	return []byte(finalMsg), nil
// }

// func MyAesDecrypt(key []byte, text []byte) ([]byte, error) {

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, errors.New(FuncRunErr("3do8awe", FuncNameErr()+" "+err.Error()))
// 	}
// 	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(string(text)))
// 	if err != nil {
// 		return nil, errors.New(FuncRunErr("mkshewjd", FuncNameErr()+" "+err.Error()))
// 	}

// 	if (len(decodedMsg) % aes.BlockSize) != 0 {
// 		return nil, errors.New(FuncRunErr("mskoeuwid", FuncNameErr()+" "+err.Error()))
// 	}

// 	iv := decodedMsg[:aes.BlockSize]
// 	msg := decodedMsg[aes.BlockSize:]

// 	cfb := cipher.NewCFBDecrypter(block, iv)
// 	cfb.XORKeyStream(msg, msg)

// 	unpadMsg, err := Unpad(msg)
// 	if err != nil {
// 		return nil, errors.New(FuncRunErr("012bsoo832d", FuncNameErr()+" "+err.Error()))
// 	}
// 	return unpadMsg, nil
// }

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

// func PkeyDecrypt(encr64 string, keyPair64 string, cliendId string) ([]byte, error) {

// 	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
// 	if err != nil {
// 		return nil, LogErr("ertjhdssw", FuncNameErr()+"Failure-1 ", err)
// 	}

// 	sKey, pKey, err := cryptobox.CryptoBoxGetSecretPublicKeyFrom(keyPair)
// 	if err != nil {
// 		return nil, LogErr("034hjalrse", FuncNameErr()+"Failure-2 ", err)
// 	}

// 	decodedBytes, err := base64.StdEncoding.DecodeString(encr64)
// 	if err != nil {
// 		return nil, LogErr("32rfww3", FuncNameErr()+"Failure-3 ", err)
// 	}

// 	decryptedBytes, boxRet := cryptobox.CryptoBoxSealOpen(decodedBytes, pKey, sKey)
// 	if boxRet != 0 {
// 		return nil, LogErr("mcnbxkajhr3eih", FuncNameErr()+"boxRet:"+strconv.Itoa(boxRet), nil)
// 	}
// 	return decryptedBytes, nil
// }

// func PkeyEncrypt(msg string, keyPair64 string, cliendId string) (string, error) {

// 	keyPair, err := base64.StdEncoding.DecodeString(keyPair64)
// 	if err != nil {
// 		return "", LogErr("2903uljslw", FuncNameErr()+"Failure-1 ", err)
// 	}

// 	_, pKey, err := cryptobox.CryptoBoxGetSecretPublicKeyFrom(keyPair)
// 	if err != nil {
// 		return "", LogErr("mcnoiaoc", FuncNameErr()+"Failure-2 ", err)
// 	}

// 	EncryptedBytes, boxRet := cryptobox.CryptoBoxSeal([]byte(msg), pKey)
// 	if boxRet != 0 {
// 		return "", LogErr("8u3h82f0", FuncNameErr()+"boxRet:"+strconv.Itoa(boxRet), nil)
// 	}
// 	EncryptedBase64 := base64.StdEncoding.EncodeToString(EncryptedBytes)
// 	return EncryptedBase64, nil
// }

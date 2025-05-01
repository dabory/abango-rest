package cryptobox

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import support "github.com/dabory/abango-rest/gosodium/support"

func CryptoBoxSeal(m []byte, pk []byte) ([]byte, int) {
	support.CheckSize(pk, CryptoBoxPublicKeyBytes(), "public key")
	c := make([]byte, len(m)+CryptoBoxSealBytes())
	exit := int(C.crypto_box_seal(
		(*C.uchar)(&c[0]),
		(*C.uchar)(&m[0]),
		(C.ulonglong)(len(m)),
		(*C.uchar)(&pk[0])))

	return c, exit
}

func CryptoBoxSealOpen(c []byte, pk []byte, sk []byte) ([]byte, int) {
	support.CheckSize(pk, CryptoBoxPublicKeyBytes(), "public key")
	support.CheckSize(sk, CryptoBoxSecretKeyBytes(), "secret key")
	m := make([]byte, len(c)-CryptoBoxSealBytes())
	exit := int(C.crypto_box_seal_open(
		(*C.uchar)(&m[0]),
		(*C.uchar)(&c[0]),
		(C.ulonglong)(len(c)),
		(*C.uchar)(&pk[0]),
		(*C.uchar)(&sk[0])))

	return m, exit
}

func CryptoBoxSealBytes() int {
	return int(C.crypto_box_sealbytes())
}

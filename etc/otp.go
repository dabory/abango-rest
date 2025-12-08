// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	OtpDigits = 10          // OTP 자릿수
	OtpPeriod = time.Minute // 1분 주기
)

func MacSecretGet() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var mac string

	for _, iface := range ifaces {
		// 루프백 제외 + MAC 존재하는 인터페이스만 사용
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 0 {
			continue
		}

		mac = iface.HardwareAddr.String()
		break
	}

	if mac == "" {
		return "", errors.New("no valid MAC address found")
	}

	// ★ MAC + BelovedPass 기반 HMAC-SHA256 Secret 생성
	h := hmac.New(sha256.New, []byte(BelovedPass))
	h.Write([]byte(mac))
	secret := hex.EncodeToString(h.Sum(nil))

	return secret, nil
}

// ============================
// 2. TOTP 생성 로직 (60초 주기)
// ============================

// 시크릿(여기서는 MAC 문자열)을 기반으로 특정 시각(now)의 TOTP 생성
func generateTOTP(secret string, digits int, period time.Duration, now time.Time) (string, error) {
	if digits <= 0 {
		return "", fmt.Errorf("digits must be > 0")
	}
	if period <= 0 {
		return "", fmt.Errorf("period must be > 0")
	}

	// 1) counter 계산 (period 단위)
	sec := int64(period.Seconds())
	counter := now.Unix() / sec

	// 2) counter를 8바이트 big-endian으로
	var counterBytes [8]byte
	binary.BigEndian.PutUint64(counterBytes[:], uint64(counter))

	// 3) HMAC-SHA1(secret, counterBytes)
	h := hmac.New(sha1.New, []byte(secret))
	if _, err := h.Write(counterBytes[:]); err != nil {
		return "", err
	}
	hash := h.Sum(nil)

	// 4) dynamic truncation (RFC 4226 스타일)
	offset := hash[len(hash)-1] & 0x0f
	code := (int(hash[offset])&0x7f)<<24 |
		(int(hash[offset+1])&0xff)<<16 |
		(int(hash[offset+2])&0xff)<<8 |
		(int(hash[offset+3]) & 0xff)

	// 5) 자릿수(digits)만큼 mod 연산
	mod := 1
	for i := 0; i < digits; i++ {
		mod *= 10
	}
	code = code % mod

	// 6) 앞에 0 채우기
	format := fmt.Sprintf("%%0%dd", digits) // 예: %06d
	return fmt.Sprintf(format, code), nil
}

// ============================
// 3. OTP Manager (메모리 저장)
// ============================

type OTPManager struct {
	mu      sync.RWMutex
	CurrOTP string // 현재 1분 구간의 OTP
	LastOTP string // 직전 1분 구간의 OTP

	Secret string        // 이 서버의 MAC 기반 시크릿
	Digits int           // 자릿수
	Period time.Duration // OTP 변경 주기
}

// 생성자: MAC 기반 Secret을 받아서 사용
func NewOTPManager(secret string, digits int, period time.Duration) *OTPManager {
	return &OTPManager{
		Secret: secret,
		Digits: digits,
		Period: period,
	}
}

// 한 번 rotate: Curr → Last, 새 Curr 생성
func (m *OTPManager) rotateOnce(now time.Time) error {
	newOTP, err := generateTOTP(m.Secret, m.Digits, m.Period, now)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.LastOTP = m.CurrOTP
	m.CurrOTP = newOTP

	log.Printf("[OTP ROTATE] LAST=%s, CURR=%s\n", m.LastOTP, m.CurrOTP)
	return nil
}

// Start: goroutine으로 1분 간격 반복
func (m *OTPManager) Start(ctx context.Context) {
	// 최초 한 번 현재 시각 기준으로 생성
	_ = m.rotateOnce(time.Now())

	// 분 경계에 맞추고 싶으면 아래처럼 조정 가능 (옵션)
	// next := time.Now().Truncate(m.Period).Add(m.Period)
	// time.Sleep(time.Until(next))

	ticker := time.NewTicker(m.Period)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("OTPManager stopped")
				return
			case now := <-ticker.C:
				if err := m.rotateOnce(now); err != nil {
					log.Println("OTPManager rotate error:", err)
				}
			}
		}
	}()
}

// Validate: CurrOTP 또는 LastOTP와 비교
func (m *OTPManager) Validate(code string) bool {
	if code == "" {
		return false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return code == m.CurrOTP || code == m.LastOTP
}

// 디버깅용: 현재 값 조회
func (m *OTPManager) DebugValues() (curr, last string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.CurrOTP, m.LastOTP
}

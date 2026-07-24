package initialization

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
)

type Config struct {
	PortName    string
	BaudRate    int
	ReadDelay   time.Duration // مدت سکوت بعد از آخرین دریافت برای تشخیص پایان پاسخ
	ReadTimeout time.Duration // timeout کلی خواندن
}

func SendAndReceive(cfg Config, req []byte) ([]byte, error) {
	mode := &serial.Mode{
		BaudRate: cfg.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(cfg.PortName, mode)
	if err != nil {
		return nil, fmt.Errorf("open serial port %s: %w", cfg.PortName, err)
	}
	defer port.Close()

	// خواندن غیرمسدودکننده/با timeout کوتاه
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 200 * time.Millisecond
	}
	if cfg.ReadDelay <= 0 {
		cfg.ReadDelay = 100 * time.Millisecond
	}

	if err := port.SetReadTimeout(cfg.ReadTimeout); err != nil {
		return nil, fmt.Errorf("set read timeout: %w", err)
	}

	// پاک کردن بافرهای قبلی اگر پشتیبانی شود
	if r, ok := port.(interface{ ResetInputBuffer() error }); ok {
		_ = r.ResetInputBuffer()
	}
	if r, ok := port.(interface{ ResetOutputBuffer() error }); ok {
		_ = r.ResetOutputBuffer()
	}

	if err := writeAll(port, req); err != nil {
		return nil, fmt.Errorf("write request: %w", err)
	}

	// کمی صبر برای پاسخ دستگاه
	time.Sleep(cfg.ReadDelay)

	resp, err := readUntilSilence(port, cfg.ReadDelay, 4*1024)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return resp, nil
}

func writeAll(port serial.Port, data []byte) error {
	for len(data) > 0 {
		n, err := port.Write(data)
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("zero-byte write")
		}
		data = data[n:]
	}
	return nil
}

func readUntilSilence(port serial.Port, silence time.Duration, maxSize int) ([]byte, error) {
	var buf bytes.Buffer
	tmp := make([]byte, 256)

	lastData := time.Now()

	for {
		n, err := port.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
			lastData = time.Now()

			if buf.Len() >= maxSize {
				return buf.Bytes(), fmt.Errorf("response too large (>%d bytes)", maxSize)
			}
		}

		// اگر timeout یا EOF نبود ولی داده نداریم، بررسی سکوت
		if err != nil {
			// بعضی درایورها timeout را به شکل error برنمی‌گردانند، پس فقط ادامه می‌دهیم
			if isTimeoutErr(err) {
				if time.Since(lastData) >= silence {
					return buf.Bytes(), nil
				}
				continue
			}
			// اگر داده گرفته‌ایم و بعد خطا آمده، بهتر است همان را برگردانیم
			if buf.Len() > 0 {
				return buf.Bytes(), nil
			}
			return nil, err
		}

		// اگر مدتی هیچ داده‌ای نیامده، پاسخ تمام شده فرض می‌کنیم
		if time.Since(lastData) >= silence && buf.Len() > 0 {
			return buf.Bytes(), nil
		}
	}
}

func readUntilDelimiter(port serial.Port, delimiter byte, maxSize int) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1)

	for len(out) < maxSize {
		n, err := port.Read(buf)
		if n > 0 {
			out = append(out, buf[:n]...)
			if out[len(out)-1] == delimiter {
				return out, nil
			}
		}
		if err != nil {
			if isTimeoutErr(err) {
				continue
			}
			if len(out) > 0 {
				return out, nil
			}
			return nil, err
		}
	}

	return out, fmt.Errorf("delimiter not found before max size")
}

func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "timeout")
}

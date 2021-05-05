package config

import (
	"strings"
	"unicode"

	"github.com/spf13/cast"
)

// parseSizeInBytes converts strings like 1GB or 12 mb into an unsigned integer number of bytes.
func parseSizeInBytes(sizeStr string) uint {
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			if lastChar > 1 {
				kbShift := 10
				mbShift := 20
				gbShift := 30

				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << kbShift
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << mbShift
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << gbShift
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size := cast.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

	return safeMul(uint(size), multiplier)
}

func safeMul(a, b uint) uint {
	c := a * b

	if a > 1 && b > 1 && c/b != a {
		return 0
	}

	return c
}

package log

// MaskPartial returns 2/3 masked varsion of a string
func MaskPartial(s string) string {
	var masked []byte
	l := len(s)
	if l <= 1 {
		return "*"
	}

	// mask 2/3 of string length
	// keep the 1/6 beginning & 1/6 end unmasked
	minMaskedIdx := l / 6
	maxMaskedIdx := l * 5 / 6
	for i, v := range s {
		if i < minMaskedIdx || i >= maxMaskedIdx {
			masked = append(masked, byte(v))
		} else {
			masked = append(masked, byte('*'))
		}
	}
	return string(masked)
}

// Mask return masked version of a string
// "aaaa" => "****"
func Mask(s string) string {
	var masked []byte
	for range s {
		masked = append(masked, byte('*'))
	}
	return string(masked)
}

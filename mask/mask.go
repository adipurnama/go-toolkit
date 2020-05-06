package mask

import (
	"bytes"
	"net/url"
)

// RedactionString represents the filtered value used in place of sensitive data in the mask package
const RedactionString = "[FILTERED]"

// IsSensitiveParam will return true if the given parameter name should be masked for sensitivity
func IsSensitiveParam(name string) bool {
	return parameterMatcher.MatchString(name)
}

// IsSensitiveHeader will return true if the given parameter name should be masked for sensitivity
func IsSensitiveHeader(name string) bool {
	return headerMatcher.MatchString(name)
}

// URL will mask the sensitive components in an URL with `[FILTERED]`.
// This list should maintain parity with the list
// Based on https://stackoverflow.com/a/52965552/474597.
func URL(originalURL string) string {
	u, err := url.Parse(originalURL)
	if err != nil {
		return "<invalid URL>"
	}

	redactionBytes := []byte(RedactionString)
	buf := bytes.NewBuffer(make([]byte, 0, len(originalURL)))

	for i, queryPart := range bytes.Split([]byte(u.RawQuery), []byte("&")) {
		if i != 0 {
			buf.WriteByte('&')
		}

		splitParam := bytes.SplitN(queryPart, []byte("="), 2)

		if len(splitParam) == 2 {
			buf.Write(splitParam[0])
			buf.WriteByte('=')

			if parameterMatcher.Match(splitParam[0]) {
				buf.Write(redactionBytes)
			} else {
				buf.Write(splitParam[1])
			}
		} else {
			buf.Write(queryPart)
		}
	}
	u.RawQuery = buf.String()
	return u.String()
}

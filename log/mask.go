package log

import (
	"bytes"
	"net/url"
	"regexp"
)

var parameterMatches = []string{
	`token$`,
	`password`,
	`secret`,
	`key$`,
	`signature`,
	`^authorization$`,
	`^certificate$`,
	`^encrypted_key$`,
	`^hook$`,
	`^import_url$`,
	`^otp_attempt$`,
	`^sentry_dsn$`,
	`^trace$`,
	`^variables$`,
	`^content$`,
	`^body$`,
	`^description$`,
	`^note$`,
	`^text$`,
	`^title$`,
}

var headerMatches = []string{
	`token$`,
	`password`,
	`secret`,
	`key$`,
	`signature`,
	`^authorization$`,
}

// parameterMatcher is precompiled for performance reasons. Keep in mind that
// `IsSensitiveParam`, `IsSensitiveHeader` and `URL` may be used in tight loops
// which may be sensitive to performance degradations.
var parameterMatcher = compileRegexpFromStrings(parameterMatches)

// headerMatcher is precompiled for performance reasons, same as `parameterMatcher`.
var headerMatcher = compileRegexpFromStrings(headerMatches)

func compileRegexpFromStrings(paramNames []string) *regexp.Regexp {
	var buffer bytes.Buffer

	buffer.WriteString("(?i)")

	for i, v := range paramNames {
		if i > 0 {
			buffer.WriteString("|")
		}

		buffer.WriteString(v)
	}

	return regexp.MustCompile(buffer.String())
}

// RedactionString represents the filtered value used in place of sensitive data in the log package.
const RedactionString = "[FILTERED]"

// IsSensitiveParam will return true if the given parameter name should be masked for sensitivity.
func IsSensitiveParam(name string) bool {
	return parameterMatcher.MatchString(name)
}

// IsSensitiveHeader will return true if the given parameter name should be masked for sensitivity.
func IsSensitiveHeader(name string) bool {
	return headerMatcher.MatchString(name)
}

// MaskURL will mask the sensitive components in an URL with `[FILTERED]`.
// This list should maintain parity with the list
// Based on https://stackoverflow.com/a/52965552/474597.
func MaskURL(originalURL string) string {
	u, err := url.Parse(originalURL)
	if err != nil {
		return "<invalid URL>"
	}

	redactionBytes := []byte(RedactionString)
	buf := bytes.NewBuffer(make([]byte, 0, len(originalURL)))

	paramSplitN := 2

	for i, queryPart := range bytes.Split([]byte(u.RawQuery), []byte("&")) {
		if i != 0 {
			buf.WriteByte('&')
		}

		splitParam := bytes.SplitN(queryPart, []byte("="), paramSplitN)

		if len(splitParam) == paramSplitN {
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

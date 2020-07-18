package mask

import (
	"bytes"
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

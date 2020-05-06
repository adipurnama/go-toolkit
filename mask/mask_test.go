package mask

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type matches uint8

const (
	whole matches = 1 << iota
	inMiddle
	atStart
	atEnd
)

const allMatches matches = whole ^ inMiddle ^ atStart ^ atEnd
const noMatches matches = 0

func TestIsSensitiveParam(t *testing.T) {
	tests := []struct {
		name string
		want matches
	}{
		{"token", whole ^ atEnd},
		{"password", allMatches},
		{"secret", allMatches},
		{"key", whole ^ atEnd},
		{"signature", allMatches},
		{"authorization", whole},
		{"certificate", whole},
		{"encrypted_key", whole ^ atEnd},
		{"hook", whole},
		{"import_url", whole},
		{"otp_attempt", whole},
		{"sentry_dsn", whole},
		{"trace", whole},
		{"variables", whole},
		{"content", whole},
		{"body", whole},
		{"description", whole},
		{"note", whole},
		{"text", whole},
		{"title", whole},
		{"gitlab", noMatches},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := func(v string) {
				// Normal case
				gotWhole := IsSensitiveParam(v)
				require.Equal(t, tt.want&whole != 0, gotWhole, "Whole param match on '%s'", v)

				gotMiddle := IsSensitiveParam("prefix" + v + "suffix")
				require.Equal(t, tt.want&inMiddle != 0, gotMiddle, "Middle match on '%s'", "prefix"+v+"suffix")

				gotStart := IsSensitiveParam(v + "suffix")
				require.Equal(t, tt.want&atStart != 0, gotStart, "Start match on '%s'", v+"suffix")

				gotEnd := IsSensitiveParam("prefix" + v)
				require.Equal(t, tt.want&atEnd != 0, gotEnd, "End match on '%s'", "prefix"+v)
			}

			check(tt.name)
			check(strings.ToUpper(tt.name))
			check(strings.ToLower(tt.name))
		})
	}
}

func TestIsSensitiveHeader(t *testing.T) {
	tests := []struct {
		name string
		want matches
	}{
		{"token", whole ^ atEnd},
		{"password", allMatches},
		{"secret", allMatches},
		{"key", whole ^ atEnd},
		{"signature", allMatches},
		{"authorization", whole},
		{"name", noMatches},
		{"gitlab", noMatches},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := func(v string) {
				gotWhole := IsSensitiveHeader(v)
				require.Equal(t, tt.want&whole != 0, gotWhole, "Whole param match on '%s'", v)

				gotMiddle := IsSensitiveHeader("prefix" + v + "suffix")
				require.Equal(t, tt.want&inMiddle != 0, gotMiddle, "Middle match on '%s'", "prefix"+v+"suffix")

				gotStart := IsSensitiveHeader(v + "suffix")
				require.Equal(t, tt.want&atStart != 0, gotStart, "Start match on '%s'", v+"suffix")

				gotEnd := IsSensitiveHeader("prefix" + v)
				require.Equal(t, tt.want&atEnd != 0, gotEnd, "End match on '%s'", "prefix"+v)
			}

			check(tt.name)
			check(strings.ToUpper(tt.name))
			check(strings.ToLower(tt.name))
		})
	}
}

func BenchmarkURL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		URL(`http://localhost:8000?token=123&something_else=92384&secret=sdmalaksjdasd&hook=123901283019238&trace=12312312312123`)
	}
}

func TestURL(t *testing.T) {
	tests := map[string]string{
		"http://localhost:8000":                                             "http://localhost:8000",
		"https://gitlab.com/":                                               "https://gitlab.com/",
		"custom://gitlab.com?secret=x":                                      "custom://gitlab.com?secret=[FILTERED]",
		"gitlab.com?secret=x":                                               "gitlab.com?secret=[FILTERED]",
		":":                                                                 "<invalid URL>",
		"http://example.com":                                                "http://example.com",
		"http://example.com?foo=1":                                          "http://example.com?foo=1",
		"http://example.com?foo=token":                                      "http://example.com?foo=token",
		"http://example.com?title=token":                                    "http://example.com?title=[FILTERED]",
		"http://example.com?authenticity_token=1":                           "http://example.com?authenticity_token=[FILTERED]",
		"http://example.com?private_token=1":                                "http://example.com?private_token=[FILTERED]",
		"http://example.com?rss_token=1":                                    "http://example.com?rss_token=[FILTERED]",
		"http://example.com?access_token=1":                                 "http://example.com?access_token=[FILTERED]",
		"http://example.com?refresh_token=1":                                "http://example.com?refresh_token=[FILTERED]",
		"http://example.com?foo&authenticity_token=blahblah&bar":            "http://example.com?foo&authenticity_token=[FILTERED]&bar",
		"http://example.com?private-token=1":                                "http://example.com?private-token=[FILTERED]",
		"http://example.com?foo&private-token=blahblah&bar":                 "http://example.com?foo&private-token=[FILTERED]&bar",
		"http://example.com?private-token=foo&authenticity_token=bar":       "http://example.com?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"https://example.com:8080?private-token=foo&authenticity_token=bar": "https://example.com:8080?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"/?private-token=foo&authenticity_token=bar":                        "/?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"?private-token=&authenticity_token=&bar":                           "?private-token=[FILTERED]&authenticity_token=[FILTERED]&bar",
		"?private-token=foo&authenticity_token=bar":                         "?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"?private_token=foo&authenticity-token=bar":                         "?private_token=[FILTERED]&authenticity-token=[FILTERED]",
		"?X-AMZ-Signature=foo":                                              "?X-AMZ-Signature=[FILTERED]",
		"?x-amz-signature=foo":                                              "?x-amz-signature=[FILTERED]",
		"?Signature=foo":                                                    "?Signature=[FILTERED]",
		"?confirmation_password=foo":                                        "?confirmation_password=[FILTERED]",
		"?pos_secret_number=foo":                                            "?pos_secret_number=[FILTERED]",
		"?sharedSecret=foo":                                                 "?sharedSecret=[FILTERED]",
		"?book_key=foo":                                                     "?book_key=[FILTERED]",
		"?certificate=foo":                                                  "?certificate=[FILTERED]",
		"?hook=foo":                                                         "?hook=[FILTERED]",
		"?import_url=foo":                                                   "?import_url=[FILTERED]",
		"?otp_attempt=foo":                                                  "?otp_attempt=[FILTERED]",
		"?sentry_dsn=foo":                                                   "?sentry_dsn=[FILTERED]",
		"?trace=foo":                                                        "?trace=[FILTERED]",
		"?variables=foo":                                                    "?variables=[FILTERED]",
		"?content=foo":                                                      "?content=[FILTERED]",
		"?content=e=mc2":                                                    "?content=[FILTERED]",
		"?formula=e=mc2":                                                    "?formula=e=mc2",
		"http://%41:8080/":                                                  "<invalid URL>",
		"https://gitlab.com?name=andrew&password=1&secret=1&key=1&signature=1&authorization=1&note=1&certificate=1&encrypted_key=1&hook=1&import_url=1&otp_attempt=1&sentry_dsn=1&trace=1&variables=1&content=1&sharedsecret=1&real=1": "https://gitlab.com?name=andrew&password=[FILTERED]&secret=[FILTERED]&key=[FILTERED]&signature=[FILTERED]&authorization=[FILTERED]&note=[FILTERED]&certificate=[FILTERED]&encrypted_key=[FILTERED]&hook=[FILTERED]&import_url=[FILTERED]&otp_attempt=[FILTERED]&sentry_dsn=[FILTERED]&trace=[FILTERED]&variables=[FILTERED]&content=[FILTERED]&sharedsecret=[FILTERED]&real=1",
	}

	for url, want := range tests {
		t.Run(url, func(t *testing.T) {
			got := URL(url)
			require.Equal(t, want, got)
		})
	}
}

package useragent_test

import (
	"testing"

	"github.com/minpeter/iplogger/pkg/useragent"
)

func TestIsCommandLine(t *testing.T) {
	// testcase stutcture
	testcase := []struct {
		ua   string
		want bool
	}{
		{"", true},
		{"curl/7.64.1", true},
		{"Wget/1.20.3", true},
		{"Mozilla/5.0 (X11; Linux x86_64; rv:72.0) Gecko/20100101 Firefox/72.0", false},
		{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) snap Chromium/79.0.3945.130 Chrome/79.0.3945.130 Safari/537.36", false},
		// Test case taken from log
		{"Expanse, a Palo Alto Networks company, searches across the global IPv4 space multiple times per day to identify customers&#39; presences on the Internet. If you would like to be excluded from our scans, please send IP addresses/domains to: scaninfo@paloaltonetworks.com", false},
	}
	for _, tc := range testcase {
		if got := useragent.IsCommandLine(tc.ua); got != tc.want {
			t.Errorf("IsCommandLine(%s) = %t, want %t", tc.ua, got, tc.want)
		}
	}
}

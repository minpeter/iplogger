package ipextractor

import (
	"net"
	"net/http"
	"strings"
)

var TrustOption = struct {
	TrustLoopback   bool
	TrustLinkLocal  bool
	TrustPrivateNet bool
	TrustCloudflare bool
	TrustIPRanges   []*net.IPNet
}{
	TrustLoopback:   true,
	TrustLinkLocal:  true,
	TrustPrivateNet: true,
	TrustCloudflare: true,
	TrustIPRanges:   []*net.IPNet{},
}

func isCloudflareRange(ip net.IP) bool {
	cloudflareIPs := []string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}

	for _, cloudflareIP := range cloudflareIPs {
		_, cloudflareRange, _ := net.ParseCIDR(cloudflareIP)
		return cloudflareRange.Contains(ip)
	}
	return false
}

func trust(ip net.IP) bool {
	if TrustOption.TrustLoopback && ip.IsLoopback() {
		return true
	}
	if TrustOption.TrustLinkLocal && ip.IsLinkLocalUnicast() {
		return true
	}
	if TrustOption.TrustPrivateNet && ip.IsPrivate() {
		return true
	}
	if TrustOption.TrustCloudflare && isCloudflareRange(ip) {
		return true
	}
	for _, trustedRange := range TrustOption.TrustIPRanges {
		return trustedRange.Contains(ip)
	}
	return false
}

func IP(r *http.Request) string {
	directIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	forwards := r.Header["X-Forwarded-For"]

	if len(forwards) == 0 {
		return directIP
	}
	ips := append(strings.Split(strings.Join(forwards, ","), ","), directIP)

	for i := len(ips) - 1; i >= 0; i-- {
		ips[i] = strings.TrimSpace(ips[i])
		ips[i] = strings.TrimPrefix(ips[i], "[")
		ips[i] = strings.TrimSuffix(ips[i], "]")
		ip := net.ParseIP(ips[i])
		if ip == nil {
			// Unable to parse IP; cannot trust entire records
			return directIP
		}
		if trust(ip) {
			return ip.String()
		}
	}

	return strings.TrimSpace(ips[0])
}

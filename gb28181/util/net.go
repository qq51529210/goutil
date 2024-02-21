package util

import "net/netip"

// IsIPV6 判断是否 ip v6
func IsIPV6(ip string) bool {
	a, _ := netip.ParseAddr(ip)
	return a.Is6()
}

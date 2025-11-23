package vpnproxy

import (
	"net"
	"strconv"
	"strings"
)

func headerDebugData(pkt []byte) []any {
	src := parseIPv4Src(pkt)
	dst := parseIPv4Dst(pkt)
	return []any{"src", src.String(), "dst", dst.String(), "header", bytesToString(pkt[:20])}
}

func parseIPv4Src(pkt []byte) net.IP {
	if len(pkt) < 20 {
		return nil
	}
	return net.IPv4(pkt[12], pkt[13], pkt[14], pkt[15])
}

func parseIPv4Dst(pkt []byte) net.IP {
	if len(pkt) < 20 {
		return nil
	}
	return net.IPv4(pkt[16], pkt[17], pkt[18], pkt[19])
}

func bytesToString(b []byte) string {
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = strconv.Itoa(int(v))
	}
	return strings.Join(parts, " ")
}

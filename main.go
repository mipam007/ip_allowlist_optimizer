package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"

	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <input_excel_file>", os.Args[0])
	}

	inputFile := os.Args[1]

	// Open Excel file
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to open Excel file: %v", err)
	}

	// Use the first sheet named "iplist"
	sheetName := "iplist"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Failed to read rows: %v", err)
	}

	// Read IPs from the sheet
	ipAddresses := []net.IP{}
	for _, row := range rows {
		if len(row) > 0 {
			ip := net.ParseIP(row[0])
			if ip != nil {
				ipAddresses = append(ipAddresses, ip)
			}
		}
	}

	// Sort IPs
	sort.Slice(ipAddresses, func(i, j int) bool {
		return compareIPs(ipAddresses[i], ipAddresses[j]) < 0
	})

	// Aggregate IPs into CIDRs
	cidrList := aggregateIPsToCIDRs(ipAddresses)

	// Print the result
	fmt.Println("Generated CIDRs:")
	for _, cidr := range cidrList {
		fmt.Printf("    \"%s\",\n", cidr)
	}
}

// compareIPs compares two IPs, returns -1 if a < b, 1 if a > b, 0 if equal
func compareIPs(a, b net.IP) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

// aggregateIPsToCIDRs aggregates a sorted list of IPs into CIDRs
func aggregateIPsToCIDRs(ips []net.IP) []string {
	cidrs := []string{}
	if len(ips) == 0 {
		return cidrs
	}

	start := ips[0]
	for i := 1; i <= len(ips); i++ {
		if i == len(ips) || !isNextIP(start, ips[i]) {
			cidr := calculateSmallestCIDR(start, ips[i-1])
			cidrs = append(cidrs, cidr)
			if i < len(ips) {
				start = ips[i]
			}
		}
	}

	return cidrs
}

// isNextIP checks if b is the next IP after a
func isNextIP(a, b net.IP) bool {
	aInt := ipToUint32(a)
	bInt := ipToUint32(b)
	return bInt == aInt+1
}

// calculateSmallestCIDR finds the smallest CIDR that includes the range from start to end
func calculateSmallestCIDR(start, end net.IP) string {
	startInt := ipToUint32(start)
	endInt := ipToUint32(end)

	for prefix := 32; prefix >= 0; prefix-- {
		mask := uint32(1<<(32-prefix)) - 1
		if startInt&^mask == startInt && startInt+mask >= endInt {
			return fmt.Sprintf("%s/%d", uint32ToIP(startInt), prefix)
		}
	}

	return fmt.Sprintf("%s/32", start.String())
}

// ipToUint32 converts an IP to a uint32
func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// uint32ToIP converts a uint32 to an IP
func uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16&0xFF), byte(n>>8&0xFF), byte(n&0xFF))
}


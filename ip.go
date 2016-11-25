package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

//The IP struct
type IP struct {
	Address     string
	CountryName string
	CommonName  string
	OrgName     string
	ServerName  string
	Delay       int
	Bandwidth   int
}

// The status of type IP
const (
	okIP = iota
	noIP
	errIP
)

//IPs []IP
type IPs []IP

//Len return the length of []IP
func (ips IPs) Len() int {
	return len(ips)
}

//Swap swap two value of []IP
func (ips IPs) Swap(i, j int) {
	ips[i], ips[j] = ips[j], ips[i]
}

//ByDelay sort by delay
type ByDelay struct {
	IPs
}

//Less return false if the first value less than the second one
func (s ByDelay) Less(i, j int) bool {
	return s.IPs[i].Delay < s.IPs[j].Delay
}

//ByBandwidth sort by bandwidth
type ByBandwidth struct {
	IPs
}

//Less return false if the first value less than the second one
func (s ByBandwidth) Less(i, j int) bool {
	return s.IPs[i].Bandwidth < s.IPs[j].Bandwidth
}

//get last ok ip
func getLastOkIP() []IP {
	m := make(map[string]IP)
	var checkedip IP
	var ips []IP
	if isFileExist(tmpOkIPFileName) {
		bytes, err := ioutil.ReadFile(tmpOkIPFileName)
		checkErr(fmt.Sprintf("read file %s error: ", tmpOkIPFileName), err, Error)
		lines := strings.Split(string(bytes), "\n")
		for _, line := range lines {
			ipInfo := strings.Split(line, " ")
			if len(ipInfo) == 6 {
				delay, err := strconv.Atoi(ipInfo[1][:len(ipInfo[1])-2])
				checkErr("delay conversion failed: ", err, Warning)
				bandwidth, err := strconv.Atoi(ipInfo[5][:len(ipInfo[5])-4])
				checkErr("bandwidth conversion failed: ", err, Warning)
				checkedip = IP{
					Address:     ipInfo[0],
					Delay:       delay,
					CommonName:  ipInfo[2],
					ServerName:  ipInfo[3],
					CountryName: ipInfo[4],
					Bandwidth:   bandwidth,
				}
				m[ipInfo[0]] = checkedip
			}
		}
	}
	for _, v := range m {
		ips = append(ips, v)
	}
	return ips
}

//get all google ip range from googleip.txt file
func getGoogleIPRange() []string {
	m := make(map[string]string)
	var ipRanges []string
	bytes, err := ioutil.ReadFile(googleIPFileName)
	checkErr(fmt.Sprintf("read file %s error: ", googleIPFileName), err, Error)

	lines := strings.Split(string(bytes), "\r\n")
	for _, line := range lines {
		if len(line) > 7 {
			m[line] = line
		}
	}
	for _, v := range m {
		ipRanges = append(ipRanges, v)
	}

	return ipRanges
}

/**
  Parse google ip range, support the following formats:
  1. xxx.xxx.xxx.xxx
  2. xxx.xxx.xxx.xxx/xx
  3. xxx.xxx.xxx.xxx-xxx.xxx.xxx.xxx
*/
func parseGoogleIPRange(ipRange string) []string {
	var ips []string
	if strings.Contains(ipRange, "/") {
		//CIDR: https://zh.wikipedia.org/wiki/%E6%97%A0%E7%B1%BB%E5%88%AB%E5%9F%9F%E9%97%B4%E8%B7%AF%E7%94%B1
		ip, ipNet, err := net.ParseCIDR(ipRange)
		checkErr(fmt.Sprintf("parse CIDR %s error: ", ipRange), err, Error)

		for iptmp := ip.Mask(ipNet.Mask); ipNet.Contains(iptmp); inc(iptmp) {
			ips = append(ips, iptmp.String())
		}
		// remove network address and broadcast address
		return ips[1 : len(ips)-1]
	} else if strings.Contains(ipRange, "-") {
		startIP := net.ParseIP(ipRange[:strings.Index(ipRange, "-")])
		endIP := net.ParseIP(ipRange[strings.Index(ipRange, "-")+1:])

		for ip := startIP; bytes.Compare(ip, endIP) <= 0; inc(ip) {
			ips = append(ips, ip.String())
		}
	} else {
		ips = append(ips, ipRange)
	}

	return ips
}

//get all google ip
func getGoogleIP() []string {
	var ips []string
	ipRanges := getGoogleIPRange()
	for _, ipRange := range ipRanges {
		ips = append(ips, parseGoogleIPRange(ipRange)...)
	}

	return ips
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

package vars

import "net"

type Service struct {
	Target   net.IP
	Port     int
	Username string
	Password string
}

type ScanResult struct {
	Server Service
	Result bool
}

type IpAddr struct {
	Ip   net.IP
	Port int
}

var Addr IpAddr
var NewIpAddr IpAddr
var Server Service

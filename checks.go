package main

import (
	//"fmt"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"time"
)

func CheckPort(host string, timeout int) (open bool) {
	conn, err := net.DialTimeout("tcp", host, time.Duration(timeout)*time.Millisecond)
	open = true
	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		open = false
		//fmt.Printf("Err: %v\n",err)
	}

	if err != nil {
		open = false
		//fmt.Printf("Err: %v\n",err)
	} else {
		conn.Close()
	}
	return open
}

func CheckIp(netinf string, ip string) (ok bool) {
	ok = false
	link, _ := netlink.LinkByName(netinf)
	addrs, _ := netlink.AddrList(link, 0)
	for _, addr := range addrs {
		if addr.String() == ip {
			if conf.Verbose {
				log.Printf("IP Found")
			}
			ok = true
			break
		}
	}
	return ok
}

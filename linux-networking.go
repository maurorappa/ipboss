package main

import (
	"github.com/vishvananda/netlink"
	"log"
	"strings"
)

func GetPrimaryIp(nic string) (ip string) {
	ip = ""
	link, err := netlink.LinkByName(nic)
	if err != nil {
		log.Printf("%q\n", err)
	}
	addr, err := netlink.AddrList(link, 2) // 2 is IPv4 from include/linux/socket.h
	// we get the first IP, our server should have only ONE IP per interface
	if verbose {
		log.Printf("Instance primary address ddresses %s:%q\n", nic, addr)
	}
	fullip := strings.Split(addr[0].String(), "/") // ugly way to gather the IP, you should use few nested struct to get *net.IPNet
	ip = fullip[0]
	return ip
}

func AddIp(network_interface string, ipaddr string) (ok bool) {
	ok = false
	netif, err := netlink.LinkByName(network_interface)
	if err != nil {
		log.Printf("\nERROR getting interface %s\n", network_interface)
	}
	addr, err := netlink.ParseAddr(ipaddr)
	if err != nil {
		log.Printf("\nERROR parse IP %s\n", ipaddr)
	}
	err = netlink.AddrAdd(netif, addr)
	if err != nil {
		log.Printf("\nERROR adding %s to %s\n", ipaddr, network_interface)
		return ok
	}
	log.Printf("\nAdding %s to %s\n", ipaddr, network_interface)
	ok = true
	return ok
}

func CheckIp(netinf string, requestedIp string) (ok bool) {
	ok = false
	link, _ := netlink.LinkByName(netinf)
	addrs, _ := netlink.AddrList(link, 0)
	for _, addr := range addrs {
		foundIP := strings.Fields(addr.String())[0]
		if foundIP == requestedIp {
			ok = true
			break
		}
	}
	return ok
}

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
	if conf.Verbose {
		log.Printf("Addresses on %s:%q\n", nic,addr)
	}
	fullip := strings.Split(addr[0].String(), "/") // ugly way to gather the IP, you should use few nested struct to get *net.IPNet
	ip=fullip[0]
	return ip
}

func AddIp(network_interface string, ipaddr string) (ok bool) {
	ok = false
	netif, err := netlink.LinkByName(network_interface)
	if err != nil {
		log.Printf("ERROR getting interface %s ", network_interface)
	}
	addr, err := netlink.ParseAddr(ipaddr)
	if err != nil {
		log.Printf("ERROR parse IP %s ", ipaddr)
	}
	err = netlink.AddrAdd(netif, addr)
	if err != nil {
		log.Printf("ERROR adding %s to %s", ipaddr, network_interface)
		return ok
	}
	log.Printf("Adding %s to %s", ipaddr, network_interface)
	ok = true
	return ok
}

func RemIp(network_interface string, ipaddr string) (ok bool) {
	ok = false
	netif, err := netlink.LinkByName(network_interface)
	if err != nil {
		log.Printf("ERROR getting interface %s ", network_interface)
	}
	addr, err := netlink.ParseAddr(ipaddr)
	if err != nil {
		log.Printf("ERROR parse IP %s ", ipaddr)
	}
	err = netlink.AddrDel(netif, addr)
	if err != nil {
		log.Printf("ERROR removing %s from %s", ipaddr, network_interface)
		return ok
	}
	log.Printf("Removing %s from %s", ipaddr, network_interface)
	ok = true
	return ok
}

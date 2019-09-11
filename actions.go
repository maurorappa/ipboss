package main

import (
	"github.com/vishvananda/netlink"
	"log"
)

func check_ip(netinf string, ip string) (ok bool) {

	link, _ := netlink.LinkByName(netinf)
	addrs, _ := netlink.AddrList(link,0)
	ok = false
	for _,addr := range addrs {
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

func add_ip(network_interface string, ip string) (ok bool) {
	ok = false
	netif, _ := netlink.LinkByName(network_interface)
	addr, _ := netlink.ParseAddr(ip + "/32")
	err := netlink.AddrAdd(netif, addr)
	if err != nil {
		log.Printf("ERROR adding %s to %s",ip,network_interface)
		return ok
	}
	log.Printf("Adding %s to %s",ip,network_interface)
	ok = true
	return ok
}

func rem_ip(network_interface string, ip string) (ok bool) {
	ok = false
	netif, _ := netlink.LinkByName(network_interface)
	addr, _ := netlink.ParseAddr(ip + "/32")
	err := netlink.AddrDel(netif, addr)
	if err != nil {
		log.Printf("ERROR removing %s from %s",ip,network_interface)
		return ok
	}
	log.Printf("Removing %s from %s",ip,network_interface)
	ok = true
	return ok
}

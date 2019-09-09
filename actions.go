package main

import (
	"github.com/vishvananda/netlink"
	"log"
)

func check_ip(netinf string, ip string) (ok bool) {

	link, _ := netlink.LinkByName(netinf)
	addrs, _ := netlink.AddrList(link,0)
	log.Printf("IPs -%v-",addrs[1])
	for _,addr := range addrs {
		if addr.String() == ip {
			log.Printf("IP Found")
		}
	}
	return ok
}

func add_ip(network_interface string, ip string) (ok bool) {
	netif, _ := netlink.LinkByName(network_interface)
	addr, _ := netlink.ParseAddr(ip + "/32")
	netlink.AddrAdd(netif, addr)
	log.Printf("Adding %s to %s",ip,network_interface)
	ok = true
	return ok
}

func rem_ip(network_interface string, ip string) (ok bool) {
	netif, _ := netlink.LinkByName(network_interface)
	addr, _ := netlink.ParseAddr(ip + "/32")
	netlink.AddrDel(netif, addr)
	log.Printf("Removing %s from %s",ip,network_interface)
	ok = true
	return ok
}

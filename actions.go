package main

import (
	"github.com/vishvananda/netlink"
)

func add_ip(network_interface string, ip string) (ok bool) {
	netif, _ := netlink.LinkByName(network_interface)
	addr, _ := netlink.ParseAddr(ip + "/32")
	netlink.AddrAdd(netif, addr)
	ok = true
	return ok
}

func rem_ip(network_interface string, ip string) (ok bool) {
	netif, _ := netlink.LinkByName(network_interface)
	addr, _ := netlink.ParseAddr(ip + "/32")
	netlink.AddrDel(netif, addr)
	ok = true
	return ok
}

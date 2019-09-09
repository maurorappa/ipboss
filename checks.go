package main

import (
	"net"
	"time"
)

func check_port(host string, timeout int) (open bool) {
	_, err := net.DialTimeout("tcp", host, time.Duration(timeout)*time.Millisecond)

	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		open = false
	}

	if err != nil {
		open = false
	}
	//conn.Close()
	return open
}

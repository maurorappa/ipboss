package main

import (
	//"fmt"
	"net"
	"time"
)

func check_port(host string, timeout int) (open bool) {
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

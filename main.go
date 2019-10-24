package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	version = "3.3"
)

var (
	verbose    bool
	verboseApi bool
	conf       *config
)

func main() {
	confPath := flag.String("c", "cfg.cfg", "Configuration file")
	service := flag.String("s", "", "service to find")
	eni := flag.Bool("e", false, "Associate an ENI to the instance running jenkins master")
	info := flag.Bool("i", false, "show info")
	flag.BoolVar(&verbose, "v", false, "Enable debug")
	flag.BoolVar(&verboseApi, "vv", false, "Enable AWS API debug")
	flag.Parse()
	conf = loadConfig(*confPath)

	if verbose || os.Getenv("VERBOSE") == "true" {
		conf.Verbose = true
	}
	//log.SetPrefix("IPboss ")
	if conf.Verbose {
		verbose = true
	}
	if verboseApi {
		verbose = true
	}
	if *info {
		fmt.Printf("Ipboss version %s\n\n", version)
	}
	if *service != "" {
		fmt.Printf("%s", findIp(*service))
	}

	if *eni {
		if !CheckIp(conf.Interface, conf.EniPrivateIp) {
			myId := findIstanceId()
			if verbose {
				fmt.Printf("This EC2 instance id is: %s\n", myId)
			}
			eniAtt, instanceId := DescEni(conf.EniId)
			if myId != instanceId {
				if !DelEni(eniAtt) {
					fmt.Printf("Cannot detach ENI")
					os.Exit(5)
				}
				time.Sleep(5 * time.Second) // we need to wait for EC2 to relinquish it
				if !AddEni(myId, conf.EniId) {
					fmt.Printf("Cannot attach ENI")
					os.Exit(7)
				}
				time.Sleep(10 * time.Second) // we need to wait for Linux to see the new interface
				fmt.Printf("waiting a bit....\n")
			} else {
				fmt.Printf("ENI already attached to this instance\n")
			}
			addRoute(conf.Interface, conf.PublicIp, conf.EniPrivateIp)
		} else {
			if verbose {
				fmt.Printf("IP for ENI already present\n")
			}
		}
	}
	_ = verbose
	_ = verboseApi
}

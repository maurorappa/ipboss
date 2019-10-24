package main

import (
	"flag"
	"fmt"
	//"os"
)

const (
	ecsCluster = "prod-exchange-jenkins"
	//netif = "wlp59s0"
	netif = "eth0"
	version = "2.2"
)

var (
	verbose bool
	verboseApi bool

)

func main() {
	service := flag.String("s", "", "service to find")
	masterIp := flag.String("m", "", "IP to be associated to the instance running jenkins master")
	flag.BoolVar(&verbose, "v", false, "Enable debug")
	flag.BoolVar(&verboseApi, "vv", false, "Enable AWS API debug")
	flag.Usage = func() {
	    fmt.Printf("Ipboss version %s\nUse `s` to specify the service you want to knwow the IP of the instance is running on\nUse `m` to assing an IP to the instance you are running this command\nUse `v` or `vv` for debugging\n\n",version)
	}
	flag.Parse()
	if verboseApi {
		verbose = true
	}
	//if os.Getenv("VERBOSE") == "true" {
	if *service != "" {
		fmt.Printf("%s", findIp(*service))
	}
	if *masterIp != "" {
		if ! CheckIp(netif, *masterIp) {
			fmt.Printf("Not %s present!", *masterIp)
			AddIp(netif,*masterIp)
			eni := FindMyEni(GetPrimaryIp(netif))
			AddIpToEni(eni,*masterIp)
		} else {
			if verbose { fmt.Printf("%s already present\n", *masterIp) }
		}
	}
	_ = verbose
	_ = verboseApi
}

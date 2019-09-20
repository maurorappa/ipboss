package main

import (
	"flag"
	"log"
	//"net"
	"os"
	"strconv"
	"time"
)

var (
	verbose bool
	eni     string
	conf    *config
)

func init() {
}

func main() {
	confPath := flag.String("c", "cfg.cfg", "Configuration file")
	verbose := flag.Bool("v", false, "Enable logging")
	flag.Parse()
	//start_time = time.Now()
	conf = loadConfig(*confPath)

	if *verbose {
		conf.Verbose = true
	}
	//log.SetPrefix("IPboss ")
	if conf.Verbose {
		log.Printf("Verbose logging is enabled")
		for servicename, _ := range conf.Rules {
			log.Printf("Service: %s monitored \n", servicename)
		}
	}
	myip := GetPrimaryIp(conf.Interface)
	log.Printf("Primary IP is: %s\n",myip)
	if conf.Aws {
		if ( os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("AWS_REGION") == "" ) {
			log.Printf("You need to specify AWS credentials and region!")
			os.Exit(17)
		}
		eni = FindMyEni(myip)
	}
	Mticker := time.NewTicker(time.Duration(conf.Poll_interval) * time.Second)
	defer Mticker.Stop()
	for range Mticker.C {
		for servicename, service := range conf.Rules {
			if conf.Verbose {
				log.Printf("Checking %s", servicename)
			}
			if service.Port != 0 {
				if conf.Verbose {
					log.Printf("is port %d open?\n", service.Port)
				}
				//check if the port is open on thge primary IP
				if CheckPort(myip+":"+strconv.Itoa(service.Port), conf.Timeout) {
					// open, add IP
					if !CheckIp(conf.Interface, service.IP+" "+conf.Interface) {
						AddIp(conf.Interface, service.IP)
						if conf.Aws {
							AddIpToEni(eni, service.IP)
						}
					}
				} else {
					// closed, remove IP
					if CheckIp(conf.Interface, service.IP+" "+conf.Interface) {
						RemIp(conf.Interface, service.IP)
						// remove from AWS, this check should be independent from the above local IP check
						if conf.Aws {
							RemIpFromEni(eni, service.IP)
						}
					}
				}
			}
			// more checks will follow
			if service.Process != "" {
				log.Printf("is %s running?\n", service.Process)
			}
		}
	}
}

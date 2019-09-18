package main

import (
	"flag"
	//"fmt"
	"log"
	//"net"
	//"os"
	"strconv"
	"time"
)

var (
	verbose bool
	//cpu_load		 int
	//iterations		 int64
	//git_info		 string
	//start_time               time.Time
	conf *config
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
	log.Printf("checking...")
	get_primary_ip(conf.Interface)
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
				//check port
				if check_port(conf.Listener+":"+strconv.Itoa(service.Port), conf.Timeout) {
					// open, add IP
					if ! check_ip(conf.Interface, service.IP+" "+conf.Interface) {
						add_ip(conf.Interface, service.IP)
					}
				} else {
					// closed, remove IP
					if check_ip(conf.Interface, service.IP+" "+conf.Interface) {
						rem_ip(conf.Interface, service.IP)
					}
				}
			}
			// more check will follow
			if service.Process != "" {
				log.Printf("is %s running?\n", service.Process)
			}
		}
	}
	log.Println("Do something...")
}

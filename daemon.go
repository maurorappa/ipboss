package main

import (
	"flag"
	//"fmt"
	"log"
	//"net"
	//"os"
	//"sync"
	"time"
)

var (
	verbose                  bool
	//cpu_load		 int
	//iterations		 int64
	//git_info		 string
	//start_time               time.Time
	conf                     *config
)

//var lock = &sync.Mutex{}

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
	Mticker := time.NewTicker(time.Duration(conf.Poll_interval) * time.Second)
	defer Mticker.Stop()
	for range Mticker.C {
		for servicename, service := range conf.Rules {
			if conf.Verbose { log.Printf("Checking %s", servicename ) }
			if service.Port != 0 {
				log.Printf("is port %d open? %v\n", service.Port, check_port("localhost:"+string(service.Port),1))
			}	
			if service.Process != "" {
				log.Printf("is %s running?\n", service.Process)
			}	
		}
	}
	log.Println("Do something...")
	//go siren_mgr()
}

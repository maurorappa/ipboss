package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var privateIPBlocks []*net.IPNet

func init() {
    for _, cidr := range []string{
        "127.0.0.0/8",    // IPv4 loopback
        "10.0.0.0/8",     // RFC1918
        "172.16.0.0/12",  // RFC1918
        "192.168.0.0/16", // RFC1918
        "::1/128",        // IPv6 loopback
        "fe80::/10",      // IPv6 link-local
        "fc00::/7",       // IPv6 unique local addr
    } {
        _, block, err := net.ParseCIDR(cidr)
        if err != nil {
            panic(fmt.Errorf("parse error on %q: %v", cidr, err))
        }
        privateIPBlocks = append(privateIPBlocks, block)
    }
}

func isPrivateIP(ip net.IP) bool {
    for _, block := range privateIPBlocks {
        if block.Contains(ip) {
            return true
        }
    }
    return false
}

func FindMyEni(myip string) (eni string) {
	eni = ""
	svc := ec2.New(session.New())
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	detail := map[string]string{}
	for idx := range result.Reservations {

		for _, inst := range result.Reservations[idx].Instances {
			if *inst.PrivateIpAddress == myip {
				//row["PrivateIP"] = *inst.PrivateIpAddress
				for _, v := range inst.Tags {
					if *v.Key == "Name" {
						detail["Name"] = *v.Value
					}
				}
				detail["InstanceId"] = *inst.InstanceId
				// we assume the EC2 instances has only one IP
				detail["Eni"] = *inst.NetworkInterfaces[0].NetworkInterfaceId
				break
			}
		}
	}
	eni = detail["Eni"]
	log.Printf("Eni: %s\n", eni)
	return eni
}

func AddIpToEni(eni string, ip string) {
	reassign := true
	theip := strings.Split(ip, "/")
	realip := theip[0]
	log.Printf("Asking AWS to assign to this instance %s\n", realip)
	sess := session.New()
	config := &aws.Config{
		Region:      aws.String("eu-west-1"),
		LogLevel:    aws.LogLevel(aws.LogDebugWithHTTPBody),
	}
	svc := ec2.New(sess, config)
	inputs := &ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(eni),
		AllowReassignment: &reassign,
		PrivateIpAddresses: []*string{
			aws.String(realip),
		},
	}
	_, errs := svc.AssignPrivateIpAddresses(inputs)
	if errs != nil {
		if aerr, ok := errs.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(errs.Error())
		}
		return
	}
	log.Println("Added")
}

func RemIpFromEni(eni string, ip string) {
	theip := strings.Split(ip, "/")
	log.Printf("Asking AWS to remove to this instance %s\n", theip[0])
	svc := ec2.New(session.New())
	inputs := &ec2.UnassignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(eni),
		PrivateIpAddresses: []*string{
			aws.String(theip[0]),
		},
	}
	_, errs := svc.UnassignPrivateIpAddresses(inputs)
	if errs != nil {
		if aerr, ok := errs.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(errs.Error())
		}
		return
	}
	log.Println("Removed")
}

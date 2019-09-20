package main

import (
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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
	fmt.Printf("Eni: %s\n", eni)
	return eni
}

func AddIpToEni(eni string, ip string) {
	reassign := true
	theip := strings.Split(ip, "/")
	fmt.Printf("Asking AWS to assign to this instance %s\n", theip[0])
	svc := ec2.New(session.New())
	inputs := &ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(eni),
		AllowReassignment: &reassign,
		PrivateIpAddresses: []*string{
			aws.String(theip[0]),
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
	fmt.Println("Added")
}

func RemIpFromEni(eni string, ip string) {
	theip := strings.Split(ip, "/")
	fmt.Printf("Asking AWS to remove to this instance %s\n", theip[0])
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
	fmt.Println("Removed")
}

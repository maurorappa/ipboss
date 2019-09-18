package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func awstest() {
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
			if *inst.PrivateIpAddress == "10.120.120.70" {
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
	fmt.Printf("%s",detail["Eni"])
	inputs := &ec2.UnassignPrivateIpAddressesInput{
	//inputs := &ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(detail["Eni"]),
		PrivateIpAddresses: []*string{
			aws.String("10.120.120.119"),
			},
		}

	//results, errs := svc.AssignPrivateIpAddresses(inputs)
	results, errs := svc.UnassignPrivateIpAddresses(inputs)
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

	fmt.Println(results)
}

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"io/ioutil"
	"log"
	"strings"
)

func findIstanceId() (id string) {
	content, err := ioutil.ReadFile("/var/lib/cloud/data/instance-id")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(content))
}

func AddEni(id string, eni string) (added bool) {
	added = false
	var config *aws.Config
	sess := session.New()
	if !verboseApi {
		config = &aws.Config{
			Region: aws.String(conf.AwsRegion),
		}
	} else {
		config = &aws.Config{
			Region:   aws.String(conf.AwsRegion),
			LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
		}
	}
	svc := ec2.New(sess, config)
	input := &ec2.AttachNetworkInterfaceInput{
		DeviceIndex:        aws.Int64(1),
		InstanceId:         aws.String(id),
		NetworkInterfaceId: aws.String(eni),
	}

	result, err := svc.AttachNetworkInterface(input)
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
		return added
	}
	added = true
	fmt.Println(result)
	return added
}

func DelEni(eni string) (added bool) {
	added = false
	var config *aws.Config
	sess := session.New()
	if !verboseApi {
		config = &aws.Config{
			Region: aws.String(conf.AwsRegion),
		}
	} else {
		config = &aws.Config{
			Region:   aws.String(conf.AwsRegion),
			LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
		}
	}
	svc := ec2.New(sess, config)
	input := &ec2.DetachNetworkInterfaceInput{
		AttachmentId: aws.String(eni),
		//Force: &true,
	}

	result, err := svc.DetachNetworkInterface(input)
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
		return added
	}
	added = true
	fmt.Println(result)
	return added
}

func DescEni(eni string) (EniAttach string, instance string) {
	var config *aws.Config
	sess := session.New()
	if !verboseApi {
		config = &aws.Config{
			Region: aws.String(conf.AwsRegion),
		}
	} else {
		config = &aws.Config{
			Region:   aws.String(conf.AwsRegion),
			LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
		}
	}
	svc := ec2.New(sess, config)
	input := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{
			aws.String(eni),
		},
	}

	result, err := svc.DescribeNetworkInterfaces(input)
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

	return *result.NetworkInterfaces[0].Attachment.AttachmentId, *result.NetworkInterfaces[0].Attachment.InstanceId
}

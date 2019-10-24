package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
)

func findIp(service string) (privateip string) {

	sess := session.New()
	var config *aws.Config
	if ! verboseApi {
		config = &aws.Config{
			Region: aws.String("eu-west-1"),
		}
	} else {
		config = &aws.Config{
			Region:   aws.String("eu-west-1"),
			LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
		}
	}
	svc := ecs.New(sess, config)
	task := &ecs.ListTasksInput{
		Cluster: aws.String(ecsCluster),
		Family:  aws.String(service),
	}
	result, err := svc.ListTasks(task)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
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

	if len(result.TaskArns) == 0 {
		fmt.Printf("Task not found!\n")
		return
	}
	arn := strings.Split(*result.TaskArns[0], "/")[1]
	if verbose {
		fmt.Printf("Task: %s\n", arn)
	}
	taskdetail := &ecs.DescribeTasksInput{
		Cluster: aws.String(ecsCluster),
		Tasks: []*string{
			aws.String(arn),
		},
	}

	resultt, err := svc.DescribeTasks(taskdetail)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
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

	if len(resultt.Tasks) == 0 {
		fmt.Printf("Container not found!\n")
		return
	}
	//fmt.Println(resultt)
	container := strings.Split(*resultt.Tasks[0].ContainerInstanceArn, "/")[1]
	if verbose {
		fmt.Printf("Container: %s\n", container)
	}
	container_detail := &ecs.DescribeContainerInstancesInput{
		Cluster: aws.String(ecsCluster),
		ContainerInstances: []*string{
			aws.String(container),
		},
	}

	resulttt, err := svc.DescribeContainerInstances(container_detail)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
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

	if len(resulttt.ContainerInstances) == 0 {
		fmt.Printf("InstanceId not found!\n")
		return
	}
	instanceid := (*resulttt.ContainerInstances[0].Ec2InstanceId)
	if verbose {
		fmt.Printf("InstanceId %s\n", instanceid)
	}
	sesss := session.New()
	svcc := ec2.New(sesss, config)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceid),
		},
	}

	resultttt, err := svcc.DescribeInstances(input)
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

	if len(resultttt.Reservations) == 0 && len(resultttt.Reservations[0].Instances) == 0 {
		fmt.Printf("Instance details not found\n")
		return
	}
	privateip = *resultttt.Reservations[0].Instances[0].PrivateIpAddress
	return privateip
}

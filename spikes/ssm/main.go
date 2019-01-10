package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		fmt.Println("exiting because $AWS_PROFILE is unset")
		os.Exit(0)
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		fmt.Println("exiting because $AWS_REGION is unset")
		os.Exit(0)
	}

	svc := ssm.New(session.Must(session.NewSession()))
	result, err := svc.GetParameters(&ssm.GetParametersInput{
		Names: []*string{
			aws.String("/foo"), aws.String("/foo/bar"), aws.String("/foo/baz"),
		},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

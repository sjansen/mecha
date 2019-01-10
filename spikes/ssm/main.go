package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	deadline := time.Now().Add(1 * time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	err := svc.GetParametersByPathPagesWithContext(ctx, &ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(1),
		Path:           aws.String("/foo"),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}, func(page *ssm.GetParametersByPathOutput, done bool) bool {
		fmt.Println(page.Parameters)
		fmt.Println("---")
		return true
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

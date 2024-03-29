package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/kenshaw/ini"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func filename() string {
	if filename := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); len(filename) != 0 {
		return filename
	}
	return config.DefaultSharedCredentialsFilename()
}

func main() {
	filename := filename()

	cfg, err := ini.LoadFile(filename)
	if err != nil {
		die(err)
	}

	ctx := context.Background()

	for _, profile := range cfg.SectionNames() {
		section := cfg.GetSection(profile)
		if kid := section.Get("aws_access_key_id"); kid == "" {
			continue
		}

		cfg, err := config.LoadDefaultConfig(
			ctx, config.WithSharedConfigProfile(profile),
		)
		if err != nil {
			die(err)
		}

		if cfg.Region == "" {
			cfg.Region = "us-east-1"
		}

		svc := iam.NewFromConfig(cfg)
		result, err := svc.ListAccessKeys(ctx, &iam.ListAccessKeysInput{})
		if err != nil {
			die(err)
		}

		fmt.Println(profile)
		for _, metadata := range result.AccessKeyMetadata {
			fmt.Printf("\taccess_key_id:\t%s\n\tusername:\t%s\n\tstatus:\t\t%s\n\tcreated:\t%s\n\n",
				aws.ToString(metadata.AccessKeyId),
				aws.ToString(metadata.UserName),
				metadata.Status,
				metadata.CreateDate,
			)
		}
	}
}

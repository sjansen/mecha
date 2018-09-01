package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/go-ini/ini"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func filename() string {
	if filename := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); len(filename) != 0 {
		return filename
	}
	return external.DefaultSharedCredentialsFilename()
}

func main() {
	filename := filename()

	cfg, err := ini.Load(filename)
	if err != nil {
		die(err)
	}

	for _, profile := range cfg.SectionStrings() {
		if profile == ini.DEFAULT_SECTION {
			continue
		}

		cfg, err := external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(profile),
		)
		if err != nil {
			die(err)
		}

		svc := iam.New(cfg)
		req := svc.ListAccessKeysRequest(&iam.ListAccessKeysInput{})
		result, err := req.Send()

		fmt.Println(profile)
		for _, metadata := range result.AccessKeyMetadata {
			fmt.Printf("\taccess_key_id:\t%s\n\tusername:\t%s\n\tstatus:\t\t%s\n\tcreated:\t%s\n\n",
				aws.StringValue(metadata.AccessKeyId),
				aws.StringValue(metadata.UserName),
				metadata.Status,
				metadata.CreateDate,
			)
		}
	}
}

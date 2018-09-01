package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/go-ini/ini"
)

func homedir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func filename() string {
	if filename := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); len(filename) != 0 {
		return filename
	}
	return filepath.Join(homedir(), ".aws", "credentials")
}

func main() {
	filename := filename()

	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, profile := range cfg.SectionStrings() {
		if profile == ini.DEFAULT_SECTION {
			continue
		}

		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials(filename, profile),
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		svc := iam.New(sess)
		result, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(profile)
		for _, metadata := range result.AccessKeyMetadata {
			fmt.Printf("\taccess_key_id:\t%s\n\tusername:\t%s\n\tstatus:\t\t%s\n\tcreated:\t%s\n\n",
				aws.StringValue(metadata.AccessKeyId),
				aws.StringValue(metadata.UserName),
				aws.StringValue(metadata.Status),
				metadata.CreateDate,
			)
		}
	}
}

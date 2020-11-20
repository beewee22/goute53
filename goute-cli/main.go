package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func main() {
	profile := flag.String("profile", "", "AWS user profile name(optional)")
	flag.Parse()

	awsConfig := &session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if *profile != "" {
		awsConfig.Profile = *profile
	}

	sess, _ := session.NewSessionWithOptions(*awsConfig)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

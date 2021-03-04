package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
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
	svc := route53.New(sess)

	result, err := svc.ListHostedZones(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	for _, hostedZone := range result.HostedZones {
		fmt.Printf("hosted zone: %s \n", *hostedZone.Name)
		id := hostedZone.Id
		recordSets, err := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{HostedZoneId: id})
		if err != nil {
			exitErrorf("record set read error %v", err.Error())
		}
		if recordSets != nil {
			for _, rs := range recordSets.ResourceRecordSets {
				fmt.Printf("\trecord set: %s\n", *rs.Name)
				for _, rr := range rs.ResourceRecords {
					fmt.Printf("\t\t%s\n", *rr.Value)
				}
			}
		} else {
			fmt.Printf("no record set with %s \n", *hostedZone.Name)
		}
	}

}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

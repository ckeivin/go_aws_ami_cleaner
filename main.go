package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	awsRegion = os.Getenv("AWS_REGION")
)

// HandleRequest handles lambda requests
func HandleRequest(ctx context.Context) (map[string]map[string]string, error) {
	fmt.Printf("Configured region is [%v]", awsRegion)

	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// client := lambda.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	// Create new EC2 client
	svc := ec2.New(sess)

	// retrive values
	input := &ec2.DescribeImagesInput{
		Owners: []*string{
			aws.String("self"),
		},
	}
	ami, err := svc.DescribeImages(input)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	// get all names and ids
	imageMap := map[string]map[string]string{}

	for _, v1 := range ami.Images {
		imageMap[*v1.ImageId] = map[string]string{}
		for _, v2 := range v1.BlockDeviceMappings {
			imageMap[*v1.ImageId][*v2.DeviceName] = *v2.Ebs.SnapshotId
		}
	}
	fmt.Printf("%+v", imageMap)
	return imageMap, nil
}
func main() {
	lambda.Start(HandleRequest)

}

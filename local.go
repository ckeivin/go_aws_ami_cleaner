package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	awsRegion = "ap-southeast-1" //os.Getenv("AWS_REGION")
	tagKey    = " name "         //os.Getenv("TAG_KEY")
	// tagValues = " windows2016-base   ;   redhat7-base " //os.Getenv("TAG_VALUES")
	tagValues = " redhat6-base   ;  redhat7-base " //os.Getenv("TAG_VALUES")
	// tagValues = " redhat7-base " //os.Getenv("TAG_VALUES")

)

// HandleRequest handles lambda requests
func main() {
	fmt.Printf("Configured region is [%v]", awsRegion)

	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// client := lambda.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	// Create new EC2 client
	svc := ec2.New(sess)
	// format environment variables
	tagKey = "tag:" + strings.TrimSpace(tagKey)
	tagValueSlice := []string{}
	for _, v := range strings.Split(tagValues, ";") {
		tagValueSlice = append(tagValueSlice, strings.TrimSpace(v)+"*")
	}
	fmt.Printf("[%v]\n", tagKey)

	for _, v := range tagValueSlice {
		fmt.Printf("[%v]\n", v)
	}

	// retrive values
	input := &ec2.DescribeImagesInput{
		Owners: []*string{
			aws.String("self"),
		},
		Filters: []*ec2.Filter{
			&ec2.Filter{
				// Name: aws.String("tag:name"),
				Name: aws.String(tagKey),

				Values: aws.StringSlice(tagValueSlice),
				// []*string{
				// 	aws.StringSlice
				// 	aws.String(tagValueSlice),
				// 	// aws.String("windows2016-base*"),
				// },
			},
		},
	}
	ami, err := svc.DescribeImages(input)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	fmt.Println("type is", reflect.TypeOf(ami))

	// get all names and ids
	imageMap := map[string]map[string]string{}

	for _, v1 := range ami.Images {
		imageMap[*v1.ImageId] = map[string]string{}
		for _, v2 := range v1.BlockDeviceMappings {
			imageMap[*v1.ImageId][*v2.DeviceName] = *v2.Ebs.SnapshotId
		}
	}
	fmt.Printf("%+v", imageMap)

}

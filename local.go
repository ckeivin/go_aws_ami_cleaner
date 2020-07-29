package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	awsRegion = "ap-southeast-1" //os.Getenv("AWS_REGION")
// tagKey    = " name "         //os.Getenv("TAG_KEY")
// tagValues = " windows2016-base   ;   redhat7-base " //os.Getenv("TAG_VALUES")
// tagValues = " redhat6-base   ;  redhat7-base " //os.Getenv("TAG_VALUES")
// tagValues = " redhat7-base " //os.Getenv("TAG_VALUES")
/*
	export AWS_REGION="ap-southeast-1"
	export TAG_KEY="name"
	export TAG_VALUES=" redhat6-base   ;  redhat7-base "
*/
)

// HandleRequest handles lambda requests
func main() {
	fmt.Printf("Configured region is [%v]\n", awsRegion)

	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// client := lambda.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	// Create new EC2 client
	svc := ec2.New(sess)
	// get tags
	tagKey, tagValues := getTags()
	// format inputs
	input := formatInput(tagKey, tagValues)
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
func getTags() (tagKey, tagValues string) {

	tagKey = os.Getenv("TAG_KEY")
	if len(tagKey) == 0 {
		fmt.Println()
		fmt.Println("'TAG_KEY' is not set. Default value of 'name' is used")
		tagKey = "name"
	} else {
		fmt.Printf("'TAG_KEY' of [%v] is set\n", tagKey)
	}
	tagValues = os.Getenv("TAG_VALUES")
	if len(tagValues) == 0 {
		fmt.Println()
		fmt.Println("'TAG_VALUES' is not set. Default value of 'windows2016-base' is used")
		tagValues = "windows2016-base"
	} else {
		fmt.Printf("'TAG_VALUES' of [%v] is set\n", tagValues)
	}
	return
}

// formatInput formats based on the tag key and values
func formatInput(tagKey, tagValues string) (input *ec2.DescribeImagesInput) {
	// format environment variables
	tagKey = "tag:" + strings.TrimSpace(tagKey)
	tagValueSlice := []string{}
	for _, v := range strings.Split(tagValues, ";") {
		tagValueSlice = append(tagValueSlice, strings.TrimSpace(v)+"*")
	}
	// debug
	fmt.Printf("tagKey: [%v]\n", tagKey)
	for i, v := range tagValueSlice {
		fmt.Printf("tagValues[%d]:[%v]\n", i, v)
	}
	// format inputs
	input = &ec2.DescribeImagesInput{
		Owners: []*string{
			aws.String("self"),
		}, Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String(tagKey),
				Values: aws.StringSlice(tagValueSlice),
			},
		},
	}
	return
}

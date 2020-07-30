package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	export AMI_AGE="15"
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
	amiAge := getAmiAge()

	// format inputs
	input := formatInput(tagKey, tagValues, amiAge)
	ami, err := svc.DescribeImages(input)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}
	// get AMI creation time
	creationDatesMap := getAmiAgeMap(ami)
	fmt.Printf("%v\n", creationDatesMap)
	// get AMI to snapshots mapping
	snapshotsMap := getSnapshotMap(ami)
	fmt.Printf("%+v\n", snapshotsMap)
	// get finalSnapshots
	finalSnapshot := getFinalSnapshotMap(amiAge, creationDatesMap, snapshotsMap)
	fmt.Printf("finalSnapshot: %+v\n", finalSnapshot)
	// summary
	fmt.Printf("Total AMIs: %v\n", len(finalSnapshot))
	totalSnapshots := 0
	for _, v1 := range finalSnapshot {
		for range v1 {
			totalSnapshots++
		}
	}
	fmt.Printf("Total Snapshots: %v\n", totalSnapshots)

	// testing zone
	// getTime := time.Now().Format(time.RFC3339)
	// getTime := time.Now()

	// fmt.Println(getTime)

	// sampleDate := "2020-02-06T08:00:39.000Z"
	// sampleDateFmt, err := time.Parse(time.RFC3339, sampleDate)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(sampleDateFmt.Format(time.RFC3339))
	// fmt.Println(reflect.TypeOf(sampleDateFmt))
	// daysDiff := getTime.Sub(sampleDateFmt)
	// fmt.Println(daysDiff.Hours() / 24)
	// fmt.Println(int(daysDiff.Hours() / 24))

}
func getFinalSnapshotMap(amiAge int, creationMap map[string]int, snapshotsMap map[string]map[string]string) (finalSnapshotMap map[string]map[string]string) {
	finalSnapshotMap = snapshotsMap
	for i, v := range creationMap {
		if v < amiAge {
			delete(finalSnapshotMap, i)
		}
	}
	return
}
func getAmiAgeMap(ami *ec2.DescribeImagesOutput) (creationMap map[string]int) {

	creationMap = map[string]int{}
	getTime := time.Now()

	for _, v1 := range ami.Images {

		creationTime, err := time.Parse(time.RFC3339, *v1.CreationDate)
		if err != nil {
			panic(err)
		}
		daysDiff := int(getTime.Sub(creationTime).Hours() / 24)
		creationMap[*v1.ImageId] = daysDiff
	}

	return
}

// getSnapshotMap gets all AMI Ids and respective snapshot Ids
func getSnapshotMap(ami *ec2.DescribeImagesOutput) (snapshotsMap map[string]map[string]string) {
	snapshotsMap = map[string]map[string]string{}

	for _, v1 := range ami.Images {
		snapshotsMap[*v1.ImageId] = map[string]string{}
		for _, v2 := range v1.BlockDeviceMappings {
			snapshotsMap[*v1.ImageId][*v2.DeviceName] = *v2.Ebs.SnapshotId
		}
	}
	return
}

func getAmiAge() (amiAge int) {

	amiAge, err := strconv.Atoi(os.Getenv("AMI_AGE"))
	if err != nil {
		fmt.Println("'AMI_AGE' is not set. Default value of '14' is used")
		amiAge = 14
	} else {

		fmt.Printf("'AMI_AGE' value of '%d' is set\n", amiAge)
	}
	println()
	return
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
func formatInput(tagKey, tagValues string, amiAge int) (input *ec2.DescribeImagesInput) {
	// format environment variables
	tagKey = "tag:" + strings.TrimSpace(tagKey)
	tagValueSlice := []string{}
	for _, v := range strings.Split(tagValues, ";") {
		tagValueSlice = append(tagValueSlice, strings.TrimSpace(v)+"*")
	}
	// debug
	fmt.Println("Filter options:")
	fmt.Printf("tagKey: [%v]\n", tagKey)
	for i, v := range tagValueSlice {
		fmt.Printf("tagValues[%d]:[%v]\n", i, v)
	}
	fmt.Printf("amiAge: [%d]\n\n", amiAge)
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

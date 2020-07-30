# go_aws_ami_cleaner
## Lambda Environment Variables

### Multiple Tags AMI Filtering
- use the "AmiTag_\<name\>" as tag key name 
- e.g. if the tagging filter requirements are 
"Name":"Web" and "Solution":"IIS" and "Solution":"Logging", then it should be set as the following:
```bash
"AmiTag_Name" : "web",
"AmiTag_Solution" : "IIS;Logging"
```

key|default value|Description
-|-|-|
AMI_AGE|14|Number of days since the creation of the AMI 
DRY_RUN|none| Whether to run the script in test mode.<br> `True` - will procceed with test mode <br> `False` - **will DELETE AMIs and Snapshots !**

## challenges
### lambda
- [x] pick up AWS session from lambda
- [x] pick up environment variables from lambda console into Go 
### code logic
- [x] check for tags and set default values
- [x] filter AMI images based on "self" tag 
- [x] filter AMI images based on tag key name and values
- [x] check age of AMI based on DAYS_OLD variable
- [ ] compare and exclude AMIs used in launch configurations
- [ ] compare and exclude AMIs used in launch templates
- [x] based on final list of AMIs, get respective snapshot IDs

### additional features
- [x] multiple tag keys and values filtering via "envSlice := os.Environ()"

## Useful links
https://docs.aws.amazon.com/lambda/latest/dg/golang-envvars.html 
https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime 
https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html

## Remember to build your handler executable for Linux!
GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main

## time format 
AWS uses ISO-8601 format 


### Lambda IAM Role Policies Required
- `AWSLambdaBasicExecutionRole`
- `CustomAMIRole` with the following policy

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "ec2:CreateTags",
            "Resource": "arn:aws:ec2:*::image/*"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeImages",
                "ec2:DeregisterImage",
                "ec2:DeleteSnapshot",
                "ec2:DescribeSnapshotAttribute",
                "autoscaling:DescribeLaunchConfigurations",
                "ec2:DescribeImageAttribute",
                "ec2:DescribeSnapshots"
            ],
            "Resource": "*"
        }
    ]
}
```
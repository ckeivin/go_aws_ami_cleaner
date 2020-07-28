# go_aws_ami_cleaner

## challenges
### lambda
- [ ] pick up AWS session from lambda
- [x] pick up environment variables from lambda console into Go 
### code logic
- [x] filter AMI images based on "self" tag 
- [ ] filter AMI images based on tag key name and values
- [ ] check age of AMI based on DAYS_OLD variable
- [ ] compare and exclude AMIs used in launch configurations
- [ ] based on final list of AMIs, get respective snapshot IDs


## Useful links
https://docs.aws.amazon.com/lambda/latest/dg/golang-envvars.html 
https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime 
https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html

## Remember to build your handler executable for Linux!
GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main
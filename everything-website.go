package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type EverythingWebsiteStackProps struct {
	awscdk.StackProps
}

type CustomCommandHooks struct{}

func (c *CustomCommandHooks) BeforeBundling(inputDir *string, outputDir *string) *[]*string {
	return &[]*string{
		jsii.String(fmt.Sprintf("cp -r %s/main/views %s/views", *inputDir, *outputDir)),
	}
}

func (c *CustomCommandHooks) AfterBundling(inputDir *string, outputDir *string) *[]*string {
	return &[]*string{}
}

func NewEverythingWebsiteStack(scope constructs.Construct, id string, props *EverythingWebsiteStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	backend := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("myGoHandler"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Entry:   jsii.String("./main"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
			CommandHooks: &CustomCommandHooks{},
		},
	})

	awsapigateway.NewLambdaRestApi(stack, jsii.String("myapi"), &awsapigateway.LambdaRestApiProps{
		Handler: backend,
	})

	return stack
}

func main() {
	defer jsii.Close()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := awscdk.NewApp(nil)

	NewEverythingWebsiteStack(app, "EverythingWebsiteStack", &EverythingWebsiteStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT_ID")),
		Region:  jsii.String("ap-southeast-1"),
	}
}

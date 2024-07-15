package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type AsmblyRegistrationEmailAutomationStackProps struct {
	awscdk.StackProps
}

func NewAsmblyRegistrationEmailAutomationStack(scope constructs.Construct, id string, props *AsmblyRegistrationEmailAutomationStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Setup the Lambda environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	lambda_env := map[string]*string{
		"MJ_APIKEY_PRIVATE": jsii.String(os.Getenv("MJ_APIKEY_PRIVATE")),
		"MJ_APIKEY_PUBLIC":  jsii.String(os.Getenv("MJ_APIKEY_PUBLIC")),
		"NEON_APIKEY":       jsii.String(os.Getenv("NEON_APIKEY")),
		"NEON_APIUSER":      jsii.String(os.Getenv("NEON_APIUSER")),
		"DEV":               jsii.String(os.Getenv("DEV")),
		"TEST_EMAIL":        jsii.String(os.Getenv("TEST_EMAIL")),
	}

	// The code that defines your stack goes here

	// Create a new container image
	dir, _ := os.Getwd()

	ecr_image := awslambda.EcrImageCode_FromAssetImage(jsii.String(filepath.Join(dir, "lambda-image")),
		&awslambda.AssetImageCodeProps{
			Platform: awsecrassets.Platform_LINUX_ARM64(),
		},
	)

	// Lambda Function
	fn := awslambda.NewFunction(stack, jsii.String("ClassRegistrationEmails"), &awslambda.FunctionProps{
		Code:         ecr_image,
		Runtime:      awslambda.Runtime_FROM_IMAGE(),
		Handler:      awslambda.Handler_FROM_IMAGE(),
		FunctionName: jsii.String("classEmailAutomation"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
		Environment:  &lambda_env,
		Architecture: awslambda.Architecture_ARM_64(),
	})

	// Lambda function URL endpoint
	fnUrl := fn.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	// Define a CloudFormation output for the function URL
	awscdk.NewCfnOutput(stack, jsii.String("ClassRegistrationEmailsUrlOutput"), &awscdk.CfnOutputProps{
		Value: fnUrl.Url(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAsmblyRegistrationEmailAutomationStack(app, "AsmblyRegistrationEmailAutomationStack", &AsmblyRegistrationEmailAutomationStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String("110104886034"),
		Region:  jsii.String("us-east-2"),
	}

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

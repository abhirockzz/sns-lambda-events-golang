package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const functionDir = "../function"

type SNSLambdaGolangStackProps struct {
	awscdk.StackProps
}

func NewSNSLambdaGolangStack(scope constructs.Construct, id string, props *SNSLambdaGolangStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("dynamodb-table"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("email"),
				Type: awsdynamodb.AttributeType_STRING},
		})

	table.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("sns-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Environment: &map[string]*string{"TABLE_NAME": table.TableName()},
			Entry:       jsii.String(functionDir),
		})

	table.GrantWriteData(function)

	snsTopic := awssns.NewTopic(stack, jsii.String("sns-topic"), nil)

	function.AddEventSource(awslambdaeventsources.NewSnsEventSource(snsTopic, nil))

	awscdk.NewCfnOutput(stack, jsii.String("sns-topic-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("sns-topic-name"),
			Value:      snsTopic.TopicName()})

	awscdk.NewCfnOutput(stack, jsii.String("dynamodb-table-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("dynamodb-table-name"),
			Value:      table.TableName()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewSNSLambdaGolangStack(app, "SNSLambdaGolangStack", &SNSLambdaGolangStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}

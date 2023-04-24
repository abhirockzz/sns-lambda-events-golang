package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var table string
var client *dynamodb.Client

func init() {
	table = os.Getenv("TABLE_NAME")
	if table == "" {
		log.Fatal("missing environment variable TABLE_NAME")
	}
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client = dynamodb.NewFromConfig(cfg)

}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {

		snsRecord := record.SNS
		fmt.Println("received message from sns", snsRecord.MessageID, "with body", snsRecord.Message)
		fmt.Println("storing message info to dynamodb table", table)

		item := make(map[string]types.AttributeValue)
		item["email"] = &types.AttributeValueMemberS{Value: snsRecord.Message}

		fmt.Println("Message Attributes:")
		for attrName, attrVal := range snsRecord.MessageAttributes {
			fmt.Println(attrName, "=", attrVal)
			//typeAndValue := make(map[string]interface{})
			attrValMap := attrVal.(map[string]interface{})

			dataType := attrValMap["Type"]

			val := attrValMap["Value"]

			switch dataType.(string) {
			case "String":
				item[attrName] = &types.AttributeValueMemberS{Value: val.(string)}
				//case "Binary":
				//item[attrName] = &types.AttributeValueMemberB{Value: val.([]byte)}
			}
		}

		_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      item,
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("item added to table")
	}
}

func main() {
	lambda.Start(handler)
}

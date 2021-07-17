package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"os"
)
type User struct {
	UserID string `json:"user_id"`
	Score int64 `json:"score"`
	CountryCode string `json:"country_code"`
}

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("localhost"),
		Endpoint: aws.String("http://localhost:8000")})
	if err != nil {
		log.Println(err)
		return
	}
	svc := dynamodb.New(sess)
	dbName := "leaderboard"

	writeUser(svc,dbName)

	params := &dynamodb.QueryInput{
		TableName:                 aws.String(dbName),
		Limit: aws.Int64(50),
		IndexName: aws.String("leaderboard_country"),
		KeyConditions: map[string]*dynamodb.Condition{
			"country_code": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("de"),
					},
				},
			},
		},
		//For sorting desc.
		ScanIndexForward: aws.Bool(false),
	}

	// Make the DynamoDB Query API call
	result, _ := svc.Query(params)

	fmt.Println(result)


}

func writeUser(svc *dynamodb.DynamoDB, dbName string)  {
	countryList := []string{
		"gb",
		"tr",
		"de",
		"us",
		"ca",
		"jp",
	}
	maxScore := 999999
	minScore := 1
	for i := 0; i < 100; i++ {
		id := uuid.New()
		rand.Shuffle(len(countryList), func(i, j int) { countryList[i], countryList[j] = countryList[j], countryList[i] })
		score := rand.Intn(maxScore-minScore) + minScore

		user := User{
			UserID: id.String(),
			Score: int64(score),
			CountryCode: countryList[0],
		}
		item, err := dynamodbattribute.MarshalMap(user)
		if err != nil {
			fmt.Println("Got error marshalling map:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item: item,
			TableName: aws.String(dbName),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}	
}
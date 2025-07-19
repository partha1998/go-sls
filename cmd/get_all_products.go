package main

import (
	"encoding/json"
	"fmt"

	"go-sls/internal/cache"
	"go-sls/internal/db"
	"go-sls/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(req)
	_ = godotenv.Load()
	_ = db.Init()
	cache.Init()

	products, err := service.GetAllProducts()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}
	data, _ := json.Marshal(products)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(data),
	}, nil
}

func main() {
	lambda.Start(handler)
}

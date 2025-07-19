package main

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"

	"go-sls/internal/cache"
	"go-sls/internal/db"
	"go-sls/internal/service"
	"go-sls/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_ = godotenv.Load()
	_ = db.Init()
	cache.Init()

	contentType := req.Headers["Content-Type"]
	mediaType, params, _ := mime.ParseMediaType(contentType)
	if mediaType != "multipart/form-data" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid content type"}, nil
	}

	body := bytes.NewReader([]byte(req.Body))
	mr := multipart.NewReader(body, params["boundary"])
	form, err := mr.ReadForm(32 << 20)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error parsing form"}, nil
	}

	files := form.File["file"]
	if len(files) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "No file uploaded"}, nil
	}

	f, _ := files[0].Open()
	defer f.Close()
	csvData, _ := io.ReadAll(f)
	products, err := utils.ParseCSV(bytes.NewReader(csvData))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	err = service.UpsertProducts(products)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Upload successful"}, nil
}

func main() {
	lambda.Start(handler)
}

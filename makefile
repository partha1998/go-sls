build:
	GOOS=linux GOARCH=amd64 go build -o bin/upload_product ./cmd/upload_product.go
	GOOS=linux GOARCH=amd64 go build -o bin/get_all_products ./cmd/get_all_products.go

deploy:
	serverless deploy

offline:
	serverless offline

local-up:
	docker compose up -d

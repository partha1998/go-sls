package utils

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"go-sls/internal/models"
)

func ParseCSV(r io.Reader) ([]models.Product, error) {
	var products []models.Product
	reader := csv.NewReader(r)
	_, _ = reader.Read() // Skip header
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		price, _ := strconv.ParseFloat(record[3], 64)
		qty, _ := strconv.Atoi(record[4])
		outOfStock := strings.ToLower(record[5]) == "true"
		p := models.Product{
			Name:       record[1],
			Image:      record[2],
			Price:      price,
			Qty:        qty,
			OutOfStock: outOfStock,
		}
		products = append(products, p)
	}
	return products, nil
}

package service

import (
	"context"
	"encoding/json"
	"log"

	"go-sls/internal/cache"
	"go-sls/internal/db"
	"go-sls/internal/models"
)

func UpsertProducts(products []models.Product) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO products (name, image, price, qty, out_of_stock)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name) DO UPDATE
		SET image = EXCLUDED.image,
			price = EXCLUDED.price,
			qty = EXCLUDED.qty,
			out_of_stock = EXCLUDED.out_of_stock
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range products {
		_, err := stmt.Exec(p.Name, p.Image, p.Price, p.Qty, p.OutOfStock)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return RefreshCache()
}

func RefreshCache() error {
	rows, err := db.DB.Query("SELECT id, name, image, price, qty, out_of_stock FROM products")
	if err != nil {
		return err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Price, &p.Qty, &p.OutOfStock); err != nil {
			return err
		}
		products = append(products, p)
	}

	jsonData, _ := json.Marshal(products)
	return cache.Client.Set(cache.Ctx, "products:all", jsonData, 0).Err()
}

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	val, err := cache.Client.Get(context.Background(), "products:all").Result()
	if err == nil {
		err := json.Unmarshal([]byte(val), &products)
		return products, err
	}
	log.Println("Cache miss. Reading from DB...")

	rows, err := db.DB.Query("SELECT id, name, image, price, qty, out_of_stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Price, &p.Qty, &p.OutOfStock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	jsonData, _ := json.Marshal(products)
	cache.Client.Set(context.Background(), "products:all", jsonData, 0)
	return products, nil
}

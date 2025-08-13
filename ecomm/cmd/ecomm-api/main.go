package main

import (
	"context"
	"log"

	"github.com/shubh-man007/ecomm/db"
	"github.com/shubh-man007/ecomm/ecomm-api/storer"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal("could not open database:", err)
	}
	defer database.Close()
	log.Println("successfully connected to database.")

	// Initialize storer
	st := storer.NewMySQLStorer(database.GetDB())

	product := &storer.Product{
		Name:         "Sample Product",
		Image:        "sample.jpg",
		Category:     "Sample Category",
		Description:  "A sample product for demonstration.",
		Rating:       5,
		NumReviews:   1,
		Price:        99.99,
		CountInStock: 10,
	}

	createdProduct, err := st.CreateProduct(context.Background(), product)
	if err != nil {
		log.Fatal("failed to create product:", err)
	}
	log.Printf("Created product: %+v\n", createdProduct)

	fetchedProduct, err := st.GetProduct(context.Background(), createdProduct.ID)
	if err != nil {
		log.Fatal("failed to fetch product:", err)
	}
	log.Printf("Fetched product: %+v\n", fetchedProduct)
}

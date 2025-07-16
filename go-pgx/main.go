package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createProductTable(db)
	product := Product{"Watch", 95.10, true}
	p_id := insertProduct(db, product)
	fmt.Printf("ID = %d\n", p_id)

	data := getProduct(db)
	fmt.Println(data)
}

func createProductTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		available BOOLEAN,
		created timestamp DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB, product Product) int {
	query := `INSERT INTO product (name, price, available)
		VALUES ($1, $2, $3) RETURNING id`

	var pk int
	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal("err")
	}

	return pk
}

func getProduct(db *sql.DB) []Product {
	data := []Product{}
	query := `SELECT name, available, price FROM product`

	rows, err := db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found")
		}
		log.Fatal()
	}

	defer rows.Close()

	var name string
	var available bool
	var price float64

	for rows.Next() {
		err := rows.Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}

	return data
}

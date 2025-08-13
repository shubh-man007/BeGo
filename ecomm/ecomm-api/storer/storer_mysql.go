package storer

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

func NewMySQLStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{
		db: db,
	}
}

func (ms *MySQLStorer) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	res, err := ms.db.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)

	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	p.ID = id

	return p, nil
}

func (ms *MySQLStorer) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product
	err := ms.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	return &p, nil
}

func (ms *MySQLStorer) ListProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product

	err := ms.db.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return products, nil
}

func (ms *MySQLStorer) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock, WHERE id=:id", p)
	if err != nil {
		return nil, fmt.Errorf("error updating product")
	}

	return p, nil
}

func (ms *MySQLStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) CreatedOrder(ctx context.Context, o *Order) (*Order, error) {
	err := ms.execTx(ctx, func(tx *sqlx.Tx) error {
		// Insert order
		orderQuery := `
			INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, created_at, updated_at)
			VALUES (:payment_method, :tax_price, :shipping_price, :total_price, :created_at, :updated_at)
		`
		res, err := tx.NamedExecContext(ctx, orderQuery, o)
		if err != nil {
			return fmt.Errorf("error inserting order: %w", err)
		}

		orderID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting order ID: %w", err)
		}
		o.ID = orderID

		// Insert order items
		itemQuery := `
			INSERT INTO order_items (name, quantity, image, price, product_id, order_id)
			VALUES (:name, :quantity, :image, :price, :product_id, :order_id)
		`

		for _, item := range o.Items {
			item.OrderID = orderID
			_, err := tx.NamedExecContext(ctx, itemQuery, item)
			if err != nil {
				return fmt.Errorf("error inserting order item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	return o, nil
}

func (ms *MySQLStorer) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := ms.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transactions: %w", rbErr)
		}
		return fmt.Errorf("error in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	return nil
}

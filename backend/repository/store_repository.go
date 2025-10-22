// repository/store.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	model "github.com/Qodarrz/fiber-app/model"
)

type StoreRepositoryInterface interface {
	// Store Items
	GetAllStoreItems(ctx context.Context, status string) ([]model.StoreItem, error)
	GetStoreItemByID(ctx context.Context, id int64) (*model.StoreItem, error)
	CreateStoreItem(ctx context.Context, item *model.StoreItem) error
	UpdateStoreItem(ctx context.Context, item *model.StoreItem) error
	UpdateStoreItemStock(ctx context.Context, id int64, newStock int) error
	DecrementStoreItemStock(ctx context.Context, id int64, quantity int) error
	IncrementStoreItemStock(ctx context.Context, id int64, quantity int) error

	// Orders
	CreateOrder(ctx context.Context, order *model.Order) error
	CreateOrderItems(ctx context.Context, orderID int64, items []model.OrderItem) error
	GetOrderByID(ctx context.Context, id int64) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]model.OrderItem, error)
	UpdateOrderStatus(ctx context.Context, id int64, status string) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(tx *sql.Tx) StoreRepositoryInterface
}

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) StoreRepositoryInterface {
	return &storeRepository{db: db}
}

func (r *storeRepository) GetAllStoreItems(ctx context.Context, status string) ([]model.StoreItem, error) {
	query := `SELECT id, name, description, price_points, stock, status, image_url, created_at 
	          FROM store_items`
	
	var args []interface{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.StoreItem
	for rows.Next() {
		var item model.StoreItem
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.PricePoints,
			&item.Stock, &item.Status, &item.ImageURL, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *storeRepository) GetStoreItemByID(ctx context.Context, id int64) (*model.StoreItem, error) {
	item := &model.StoreItem{}
	query := `SELECT id, name, description, price_points, stock, status, image_url, created_at 
	          FROM store_items WHERE id = $1`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.Name, &item.Description, &item.PricePoints,
		&item.Stock, &item.Status, &item.ImageURL, &item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *storeRepository) CreateStoreItem(ctx context.Context, item *model.StoreItem) error {
	query := `INSERT INTO store_items (name, description, price_points, stock, status, image_url, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	
	return r.db.QueryRowContext(ctx, query,
		item.Name, item.Description, item.PricePoints, item.Stock,
		item.Status, item.ImageURL, item.CreatedAt,
	).Scan(&item.ID)
}

func (r *storeRepository) UpdateStoreItem(ctx context.Context, item *model.StoreItem) error {
	query := `UPDATE store_items 
	          SET name = $1, description = $2, price_points = $3, stock = $4, status = $5, image_url = $6 
	          WHERE id = $7`
	
	_, err := r.db.ExecContext(ctx, query,
		item.Name, item.Description, item.PricePoints, item.Stock,
		item.Status, item.ImageURL, item.ID,
	)
	return err
}

func (r *storeRepository) UpdateStoreItemStock(ctx context.Context, id int64, newStock int) error {
	query := `UPDATE store_items SET stock = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newStock, id)
	return err
}

func (r *storeRepository) DecrementStoreItemStock(ctx context.Context, id int64, quantity int) error {
	query := `UPDATE store_items SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
	result, err := r.db.ExecContext(ctx, query, quantity, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("insufficient stock or item not found")
	}
	return nil
}

func (r *storeRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	query := `INSERT INTO orders (user_id, total_points, status, created_at)
	          VALUES ($1, $2, $3, $4) RETURNING id`
	
	return r.db.QueryRowContext(ctx, query,
		order.UserID, order.TotalPoints, order.Status, order.CreatedAt,
	).Scan(&order.ID)
}

func (r *storeRepository) CreateOrderItems(ctx context.Context, orderID int64, items []model.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO order_items (order_id, item_id, qty, price_each_points, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	
	for _, item := range items {
		_, err := tx.ExecContext(ctx, query,
			orderID, item.ItemID, item.Qty, item.PriceEachPoints, item.CreatedAt,
		)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (r *storeRepository) GetOrderByID(ctx context.Context, id int64) (*model.Order, error) {
	order := &model.Order{}
	query := `SELECT id, user_id, total_points, status, created_at FROM orders WHERE id = $1`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.UserID, &order.TotalPoints, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}

func (r *storeRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error) {
	query := `SELECT id, user_id, total_points, status, created_at 
	          FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.TotalPoints, &order.Status, &order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *storeRepository) GetOrderItems(ctx context.Context, orderID int64) ([]model.OrderItem, error) {
	query := `SELECT oi.id, oi.order_id, oi.item_id, oi.qty, oi.price_each_points, oi.created_at
	          FROM order_items oi
	          WHERE oi.order_id = $1`
	
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ItemID,
			&item.Qty, &item.PriceEachPoints, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *storeRepository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *storeRepository) IncrementStoreItemStock(ctx context.Context, id int64, quantity int) error {
	query := `UPDATE store_items SET stock = stock + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, quantity, id)
	return err
}

func (r *storeRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *storeRepository) WithTx(tx *sql.Tx) StoreRepositoryInterface {
	return &storeRepository{db: r.db} // Note: This would need proper implementation for transaction support
}
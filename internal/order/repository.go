package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *Order) error {
	const op = "repository.order.SaveOrder"

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	// Сохраняем основной заказ
	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service,
			shardkey, sm_id, date_created, oof_shard
		) VALUES (
			:order_uid, :track_number, :entry, :locale,
			:internal_signature, :customer_id, :delivery_service,
			:shardkey, :sm_id, :date_created, :oof_shard
		)`, order)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrOrderExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем доставку
	_, err = tx.ExecContext(ctx, `
        INSERT INTO deliveries (
            order_uid, name, phone, zip, city, 
            address, region, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
        order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
        order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем оплату
	_, err = tx.ExecContext(ctx, `
	INSERT INTO payments (
		order_uid, transaction, request_id, currency,
		provider, amount, payment_dt, bank, 
		delivery_cost, goods_total, custom_fee
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
	order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
	order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
	order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем товары
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id,
			brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID,
		item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID,
		item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrderByUID(ctx context.Context, uid string) (*Order, error) {
	const op = "repository.order.GetOrderByUID"

	var order Order
	err := r.db.GetContext(ctx, &order, `
		SELECT * FROM orders WHERE order_uid = $1`, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Загружаем доставку
	err = r.db.GetContext(ctx, &order.Delivery, `
		SELECT * FROM deliveries WHERE order_uid = $1`, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Загружаем оплату
	err = r.db.GetContext(ctx, &order.Payment, `
		SELECT * FROM payments WHERE order_uid = $1`, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Загружаем товары
	err = r.db.SelectContext(ctx, &order.Items, `
		SELECT * FROM items WHERE order_uid = $1`, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &order, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]Order, error) {
	const op = "repository.order.GetAllOrders"

	var orders []Order
	err := r.db.SelectContext(ctx, &orders, `SELECT * FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Для каждого заказа подгружаем связанные данные
	for i := range orders {
		uid := orders[i].OrderUID
		
		err = r.db.GetContext(ctx, &orders[i].Delivery, `
			SELECT * FROM deliveries WHERE order_uid = $1`, uid)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = r.db.GetContext(ctx, &orders[i].Payment, `
			SELECT * FROM payments WHERE order_uid = $1`, uid)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = r.db.SelectContext(ctx, &orders[i].Items, `
			SELECT * FROM items WHERE order_uid = $1`, uid)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return orders, nil
}

var ErrOrderExists = errors.New("order already exists")
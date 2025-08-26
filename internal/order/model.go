package order

import (
	"time"
	_ "encoding/json"
)


type Order struct {
	OrderUID          string    `json:"order_uid" db:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" db:"track_number" validate:"required"`
	Entry             string    `json:"entry" db:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" db:"-" validate:"required"`
	Payment           Payment   `json:"payment" db:"-" validate:"required"`
	Items             []Item    `json:"items" db:"-" validate:"required,dive,required"`
	Locale            string    `json:"locale" db:"locale" validate:"required,oneof=en ru"`
	CustomerID        string    `json:"customer_id" db:"customer_id" validate:"required"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" db:"shardkey" validate:"required"`
	SMID              int       `json:"sm_id" db:"sm_id" validate:"required,min=1"`
	DateCreated       time.Time `json:"date_created" db:"date_created" validate:"required"`
	OOFShard          string    `json:"oof_shard" db:"oof_shard" validate:"required"`
}

type Delivery struct {
	ID      int    `json:"-" db:"id"`
	OrderUID string `json:"-" db:"order_uid"`
	Name    string `json:"name" db:"name" validate:"required"`
	Phone   string `json:"phone" db:"phone" validate:"required,e164"`
	Zip     string `json:"zip" db:"zip" validate:"required"`
	City    string `json:"city" db:"city" validate:"required"`
	Address string `json:"address" db:"address" validate:"required"`
	Region  string `json:"region" db:"region" validate:"required"`
	Email   string `json:"email" db:"email" validate:"required,email"`
}

type Payment struct {
	ID           int    `json:"-" db:"id"`
	OrderUID     string `json:"-" db:"order_uid"`
	Transaction  string `json:"transaction" db:"transaction" validate:"required"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency" validate:"required"`
	Provider     string `json:"provider" db:"provider" validate:"required"`
	Amount       int    `json:"amount" db:"amount" validate:"required,min=0"`
	PaymentDT    int64  `json:"payment_dt" db:"payment_dt" validate:"required,min=0"`
	Bank         string `json:"bank" db:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" validate:"min=0"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" validate:"min=0"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee" validate:"min=0"`
}

type Item struct {
	ID         int    `json:"-" db:"id"`
	OrderUID   string `json:"-" db:"order_uid"`
	ChrtID     int    `json:"chrt_id" db:"chrt_id" validate:"required,min=1"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required"`
	Price      int    `json:"price" db:"price" validate:"required,min=0"`
	RID        string `json:"rid" db:"rid" validate:"required"`
	Name       string `json:"name" db:"name" validate:"required"`
	Sale       int    `json:"sale" db:"sale" validate:"min=0,max=100"`
	Size       string `json:"size" db:"size" validate:"required"`
	TotalPrice int    `json:"total_price" db:"total_price" validate:"min=0"`
	NmID       int    `json:"nm_id" db:"nm_id" validate:"required,min=1"`
	Brand      string `json:"brand" db:"brand" validate:"required"`
	Status     int    `json:"status" db:"status" validate:"required,min=0"`
}
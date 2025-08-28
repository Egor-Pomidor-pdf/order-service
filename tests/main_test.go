package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq" 

)

const (
	orderUID              = "b563feb7b2b84b6test1111"
	kafkaBootstrapServers = "localhost:9093"
	kafkaTopic            = "orders"
	orderMessage          = `{
   "order_uid": "b563feb7b2b84b6test1111",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`
)

func Test_OrderFlow(t *testing.T) {

	db := connectToDB(t)
	defer db.Close()

	// 1) Отправляем сообщение в Kafka
	sendOrderToKafka(t, kafkaBootstrapServers, kafkaTopic, orderMessage)

	// 2) Ждем секунду чтобы сообщение обработалось
    time.Sleep(1 * time.Second)
    
    // 3) Проверяем что заказ есть в БД
	var orderUID_2 string
    err := db.QueryRow(
        "SELECT order_uid FROM orders WHERE order_uid = $1",
        orderUID,
    ).Scan(&orderUID_2)

    require.NoError(t, err)
    require.Equal(t, "b563feb7b2b84b6test1111", orderUID_2)

	

}

func connectToDB(t *testing.T) *sql.DB {
	host := "localhost"
    port := "5430"  
    user := "order_user"
    password := "order_password"
    dbname := "order_db"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	err = db.Ping()
	require.NoError(t, err)
	return db
}

func sendOrderToKafka(t *testing.T, bootstrapServers, topic, message string) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(bootstrapServers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := w.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		require.NoError(t, err)
	} else {
		t.Log("Message sent to Kafka successfully")
	}

}

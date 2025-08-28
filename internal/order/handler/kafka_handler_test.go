package handler

import (
	"encoding/json"
	"testing"

	"github.com/Egor-Pomidor-pdf/order-service/internal/order"
	mock_service "github.com/Egor-Pomidor-pdf/order-service/internal/order/service/mocks"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleMessage(t *testing.T) {
type mockBehavior func(s *mock_service.MockOrderServiceInterface, order order.Order)
	tests := []struct {
		name         string
		message []byte
     	 offset   int64
		mockBehavior mockBehavior
		expectedErr  error
	}{
		{
			name:    "Valid Message",
			message: []byte(`{
   "order_uid": "b563feb7b2b84b6test10",
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
}`),
			offset:  1,
			mockBehavior: func(s *mock_service.MockOrderServiceInterface, order order.Order) {
				s.EXPECT().ProcessOrder(gomock.Any(), order).Return(nil)
			},
			expectedErr: nil,

		},
		// other test cases...
		
	}

	for _, test := range tests {
		t.Run(test.name, func(t * testing.T) {
			// init controller
			c := gomock.NewController(t)
			defer c.Finish()



			// init mock
			mockOrderService := mock_service.NewMockOrderServiceInterface(c)
			var expectedOrder order.Order
			json.Unmarshal(test.message, &expectedOrder)

			test.mockBehavior(mockOrderService, expectedOrder)

			// init handler
			handler := &OrderHandler{
				service: mockOrderService,
			}

			errKafka := handler.HandleMessage(test.message, kafka.Offset(test.offset))
			assert.Equal(t, test.expectedErr, errKafka)

			
		})
	}
}


package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Egor-Pomidor-pdf/order-service/internal/order"
	mock_service "github.com/Egor-Pomidor-pdf/order-service/internal/order/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)



func TestHandler_GetOrderHandler(t *testing.T) {
// init test Table
type mockBehavior func(s *mock_service.MockOrderServiceInterface, order_uid string)
	tests := []struct {
		name                 string
		order_uid            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK", 
			order_uid: "b563feb7b2b84b6test10",
			mockBehavior: func(s *mock_service.MockOrderServiceInterface, order_uid string) {
				s.EXPECT().GetOrder(gomock.Any(), order_uid).Return(&order.Order{
					OrderUID: "b563feb7b2b84b6test10",
					TrackNumber: "WBILMTESTTRACK",
					Entry: "WBIL",
					Delivery: order.Delivery{
						Name: "Test Testov",
						Phone: "+9720000000",
						Zip: "2639809",
						City: "Kiryat Mozkin",
						Address: "Ploshad Mira 15",
						Region: "Kraiot",
						Email: "test@gmail.com",
					},
					Payment: order.Payment{
						Transaction: "b563feb7b2b84b6test",
						RequestID: "",
						Currency: "USD",
						Provider: "wbpay",
						Amount: 1817,
						PaymentDT: 1637907727,
						Bank: "alpha",
						DeliveryCost: 1500,
						GoodsTotal: 317,
						CustomFee: 0,
					},
					Items: []order.Item{
						{
							ChrtID: 9934930,
							TrackNumber: "WBILMTESTTRACK",
							Price: 453,
							RID: "ab4219087a764ae0btest",
							Name: "Mascaras",
							Sale: 30,
							Size: "0",
							TotalPrice: 317,
							NmID: 2389212,
							Brand: "Vivienne Sabo",
							Status: 202,
						},
					},
					Locale: "en",
					InternalSignature: "",
					CustomerID: "test",
					DeliveryService: "meest",
					ShardKey: "9",
					SMID: 99,	
					DateCreated: time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
					OOFShard: "1",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{
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
}`,
		},
		{
			name: "NotFound",
			order_uid: "nonexistent_order",
			mockBehavior: func(s *mock_service.MockOrderServiceInterface, order_uid string) {
				s.EXPECT().GetOrder(gomock.Any(), order_uid).Return(nil, nil)
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"error": "Order not found"}`,
		},
		{
			name: "Service Error",
			order_uid: "b563feb7b2b84b6test10",
			mockBehavior: func(s *mock_service.MockOrderServiceInterface, order_uid string) {
				s.EXPECT().GetOrder(gomock.Any(), order_uid).Return(nil, errors.New("Service Error"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"error": "Internal server error"}`,
		},
		
	}

	for _, test := range tests {
		t.Run(test.name, func(t * testing.T) {
			// init controller
			c := gomock.NewController(t)
			defer c.Finish()



			// init mock
			mockOrderService := mock_service.NewMockOrderServiceInterface(c)
			test.mockBehavior(mockOrderService, test.order_uid)

			// init handler
			handler := &OrderHandler{
				service: mockOrderService,
			}

			// create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/orders/" + test.order_uid, nil)
			slog.Info("Url_Id", "id", test.order_uid)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("order_uid", test.order_uid)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			handler.GetOrderHandler(w, r)

			// check response
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())

		})
	}
}
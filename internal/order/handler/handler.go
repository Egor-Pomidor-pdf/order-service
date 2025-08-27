package handler

import "github.com/Egor-Pomidor-pdf/order-service/internal/order/service"

// type OrderHTTPHandler struct {
//     service *service.OrderService
// }

// func NewOrderHTTPHandler(service *service.OrderService) *OrderHTTPHandler {
//     return &OrderHTTPHandler{service: service}
// }

type OrderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(service service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{service: service}
}
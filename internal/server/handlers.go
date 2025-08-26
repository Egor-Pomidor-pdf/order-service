package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Egor-Pomidor-pdf/order-service/internal/order"
	"github.com/go-chi/chi/v5"
)

type OrderHTTPHandler struct {
    service *order.OrderService
}

func NewOrderHTTPHandler(service *order.OrderService) *OrderHTTPHandler {
    return &OrderHTTPHandler{service: service}
}

func (h *OrderHTTPHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	order_uid := chi.URLParam(r, "order_uid")
	order, err := h.service.GetOrder(r.Context(), order_uid)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error("failed to get order", "error", err, "order_uid", order_uid) 
		w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	if order == nil {
		slog.Info("order not found", "order_uid", order_uid)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return	
	}
	
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(order) 

}

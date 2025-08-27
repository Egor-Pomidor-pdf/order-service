package server

import (
	"net/http"
	"time"

	"github.com/Egor-Pomidor-pdf/order-service/internal/config"
	"github.com/Egor-Pomidor-pdf/order-service/internal/order/handler"
	"github.com/Egor-Pomidor-pdf/order-service/internal/order/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(cfg config.ServerConfig, orderService *service.OrderService) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	orderHandler := handler.NewOrderHandler(orderService)

	router.Get("/order/{order_uid}", orderHandler.GetOrderHandler)

	router.Handle("/*", http.FileServer(http.Dir("./web")))

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second, 
	}
		return srv
		
}	





 


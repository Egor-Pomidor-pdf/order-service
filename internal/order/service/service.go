package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/Egor-Pomidor-pdf/order-service/internal/order"
	"github.com/Egor-Pomidor-pdf/order-service/internal/order/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type OrderServiceInterface interface {
	ProcessOrder(ctx context.Context, order order.Order) error
	GetOrder(ctx context.Context, uid string) (*order.Order, error)
}

type OrderService struct {
	repo  repository.OrderRepository
	cache map[string]order.Order
	mu    sync.RWMutex
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	service := &OrderService{
		repo:  *repo,
		cache: make(map[string]order.Order),
	}
	service.restoreCache(context.Background())
	return service
}

func (s *OrderService) restoreCache(ctx context.Context) {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, order := range orders {
		s.cache[order.OrderUID] = order
	}
	slog.Info("Restore cache")
}


func (s *OrderService) ProcessOrder(ctx context.Context, order order.Order) error {
	if err := s.repo.SaveOrder(ctx, &order); err != nil {
		return err
	}

	s.mu.Lock()
	s.cache[order.OrderUID] = order
	s.mu.Unlock()

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, uid string) (*order.Order, error) {
	// Проверяем кэш
	s.mu.RLock()
	cachedOrder, exists := s.cache[uid]
	s.mu.RUnlock()

	if exists {
		slog.Info(" Order from cache\n","uid", uid)
		return &cachedOrder, nil
	}

	// Если нет в кэше, ищем в БД
	order, err := s.repo.GetOrderByUID(ctx, uid)
	slog.Info(" Order from bd\n","uid", uid)
	if err != nil {
		return nil, errors.New("Service Error")
	}

	// Обновляем кэш
	if order != nil {
		s.mu.Lock()
		s.cache[order.OrderUID] = *order
		s.mu.Unlock()
	}

	return order, nil
}
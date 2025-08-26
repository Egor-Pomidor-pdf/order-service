package order

import (
	"context"
	"sync"
)

type OrderService struct {
	repo  OrderRepository
	cache map[string]Order
	mu    sync.RWMutex
}

func NewOrderService(repo OrderRepository) *OrderService {
	service := &OrderService{
		repo:  repo,
		cache: make(map[string]Order),
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
}


func (s *OrderService) ProcessOrder(ctx context.Context, order Order) error {
	if err := s.repo.SaveOrder(ctx, &order); err != nil {
		return err
	}

	s.mu.Lock()
	s.cache[order.OrderUID] = order
	s.mu.Unlock()

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, uid string) (*Order, error) {
	// Проверяем кэш
	s.mu.RLock()
	cachedOrder, exists := s.cache[uid]
	s.mu.RUnlock()

	if exists {
		return &cachedOrder, nil
	}

	// Если нет в кэше, ищем в БД
	order, err := s.repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Обновляем кэш
	if order != nil {
		s.mu.Lock()
		s.cache[order.OrderUID] = *order
		s.mu.Unlock()
	}

	return order, nil
}
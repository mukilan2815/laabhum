package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/google/uuid"
)

type Order struct {
    ID        string
    Symbol    string
    Status    models.OrderStatus
    Strategy  models.TradeStrategy
    CreatedAt time.Time
trades           map[string][]models.Trade
}

func (r *InMemoryOrderRepository) GetTrades(parentID string) ([]models.Trade, error) {
	trades, ok := r.trades[parentID]
	if !ok {
		return nil, nil // or return an error if needed
	}
	return trades, nil
}

type OrderRepository interface {
    GetOrder(id string) (*models.Order, error)
    UpdateOrder(order models.Order) error
    DeleteOrder(id string) error
    CreatePosition(position models.Position) error
    GetPosition(id string) (*models.Position, error)
    UpdatePosition(position models.Position) error
    GetOpenPositions() ([]models.Position, error)
    ClosePosition(id string) error
    CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error)
    SaveOrder(order map[string]interface{}) error
    GetTrades(parentID string) ([]models.Trade, error)
    SaveMarketCondition(condition models.MarketCondition) error
    GetLatestMarketCondition(symbol string) (*models.MarketCondition, error)
        GetOrders(filter OrderFilter) ([]models.Order, error)
    CreateOrder(order models.Order) (models.Order, error)
    ExecuteChildOrder(orderID string) error // Add this method signature

    UpdateOrderStatus(id string, status models.OrderStatus) error // Add this method to the interface
}
func (r *InMemoryOrderRepository) ExecuteChildOrder(orderID string) error {
    order, exists := r.orders[orderID]
    if !exists {
        return fmt.Errorf("order not found")
    }
    // Implement logic for executing child order here
    // For example, update order status to "executed"
    order.Status = "executed"
    r.orders[orderID] = order
    return nil
}

type OrderFilter struct {
    Symbol   string
    Status   models.OrderStatus
    Strategy models.TradeStrategy
    FromDate time.Time
    ToDate   time.Time
}

type InMemoryOrderRepository struct {
    orders           map[string]*models.Order
    positions        map[string]*models.Position
    marketConditions map[string]*models.MarketCondition
    trades           map[string][]models.Trade
    mutex            sync.RWMutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{
        orders:           make(map[string]*models.Order),
        positions:        make(map[string]*models.Position),
        marketConditions: make(map[string]*models.MarketCondition),
        trades:           make(map[string][]models.Trade),
    }
}

func (r *InMemoryOrderRepository) SaveOrder(order map[string]interface{}) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    id, ok := order["ID"].(string)
    if !ok || id == "" {
        return errors.New("invalid order ID")
    }

    // Assuming order is a map with string keys and interface{} values
    // You might need to convert this map to your Order struct
    // This is just a placeholder implementation
    r.orders[id] = &models.Order{
        ID:        id,
        Symbol:    order["Symbol"].(string),
        Status:    order["Status"].(models.OrderStatus),
        Strategy:  order["Strategy"].(models.TradeStrategy),
        CreatedAt: time.Now().Unix(),
    }

    return nil
}
func (repo *InMemoryOrderRepository) CreateOrder(order models.Order) (models.Order, error) {
    repo.mutex.Lock()
    defer repo.mutex.Unlock()

    if order.ID == "" {
        order.ID = uuid.New().String()
    }
    order.CreatedAt = time.Now().Unix()
    repo.orders[order.ID] = &order // Keep the order in the map as a pointer, but return as a value

    return order, nil // Return the order as a value, not a pointer
}
func (repo *InMemoryOrderRepository) CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error) {

    // Implement the method to satisfy the OrderRepository interface

    return &order, nil

}


func (r *InMemoryOrderRepository) GetOrder(id string) (*models.Order, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    order, exists := r.orders[id]
    if !exists {
        return nil, errors.New("order not found")
    }
    return order, nil
}

func (r *InMemoryOrderRepository) UpdateOrder(order models.Order) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.orders[order.ID]; !exists {
        return errors.New("order not found")
    }
    r.orders[order.ID] = &order
    return nil
}

func (r *InMemoryOrderRepository) GetOrders(filter OrderFilter) ([]models.Order, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    var orders []models.Order
    for _, order := range r.orders {
        if r.orderMatchesFilter(order, filter) {
            orders = append(orders, *order)
        }
    }
    return orders, nil
}

func (r *InMemoryOrderRepository) orderMatchesFilter(order *models.Order, filter OrderFilter) bool {
    if filter.Symbol != "" && order.Symbol != filter.Symbol {
        return false
    }
    if filter.Status != "" && order.Status != filter.Status {
        return false
    }
    if filter.Strategy != "" && order.Strategy != filter.Strategy {
        return false
    }
    if !filter.FromDate.IsZero() && time.Unix(order.CreatedAt, 0).Before(filter.FromDate) {
        return false
    }
    if !filter.ToDate.IsZero() && time.Unix(order.CreatedAt, 0).After(filter.ToDate) {
        return false
    }
    return true
}

func (r *InMemoryOrderRepository) DeleteOrder(id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.orders[id]; !exists {
        return errors.New("order not found")
    }
    delete(r.orders, id)
    return nil
}

func (r *InMemoryOrderRepository) CreatePosition(position models.Position) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if position.ID == "" {
        position.ID = uuid.New().String()
    }
    position.OpenedAt = time.Now()
    position.LastUpdatedAt = time.Now()
    r.positions[position.ID] = &position
    return nil
}

func (r *InMemoryOrderRepository) GetPosition(id string) (*models.Position, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    position, exists := r.positions[id]
    if !exists {
        return nil, errors.New("position not found")
    }
    return position, nil
}

func (r *InMemoryOrderRepository) UpdatePosition(position models.Position) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.positions[position.ID]; !exists {
        return errors.New("position not found")
    }
    position.LastUpdatedAt = time.Now()
    r.positions[position.ID] = &position
    return nil
}

func (r *InMemoryOrderRepository) GetOpenPositions() ([]models.Position, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    var positions []models.Position
    for _, position := range r.positions {
        positions = append(positions, *position)
    }
    return positions, nil
}

func (r *InMemoryOrderRepository) ClosePosition(id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    position, exists := r.positions[id]
    if !exists {
        return errors.New("position not found")
    }
    delete(r.positions, id)
    
    // You might want to create a closed position history here
    fmt.Printf("Position closed: %+v\n", position)
    
    return nil
}

func (r *InMemoryOrderRepository) SaveMarketCondition(condition models.MarketCondition) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    r.marketConditions[condition.Symbol] = &condition
    return nil
}

func (r *InMemoryOrderRepository) GetLatestMarketCondition(symbol string) (*models.MarketCondition, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    condition, exists := r.marketConditions[symbol]
    if !exists {
        return nil, errors.New("market condition not found for symbol")
    }
    return condition, nil
}

func (r *InMemoryOrderRepository) UpdateOrderStatus(id string, status models.OrderStatus) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    order, exists := r.orders[id]
    if !exists {
        return errors.New("order not found")
    }
    order.Status = status
    return nil
}
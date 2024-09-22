package service

import (
	"errors"
	"time"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/google/uuid"
)

type OMSService struct {
    repo repository.OrderRepository
}



func NewOMSService(repo repository.OrderRepository) *OMSService {
    return &OMSService{repo: repo}
}

type InMemoryOrderRepository struct {
    orders        map[string]models.Order
    scalperOrders map[string]models.ScalperOrder
    trades        map[string][]models.Trade
    positions     map[string]models.Position
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{
        orders:        make(map[string]models.Order),
        scalperOrders: make(map[string]models.ScalperOrder),
        trades:        make(map[string][]models.Trade),
        positions:     make(map[string]models.Position),
    }
}

func (r *InMemoryOrderRepository) CreateOrder(order models.Order) (models.Order, error) {
    r.orders[order.ID] = order
    return order, nil
}

func (r *InMemoryOrderRepository) CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error) {
    r.scalperOrders[order.ID] = order
    return &order, nil
}

func (r *InMemoryOrderRepository) ExecuteChildOrder(parentID, childID string) error {
    // Implement the method to satisfy the OrderRepository interface
    return nil
}
func (r *InMemoryOrderRepository) GetTrades(parentID string) ([]models.Trade, error) {
    trades, ok := r.trades[parentID]
    if !ok {
        return nil, errors.New("trades not found for the given parentID")
    }
    return trades, nil
}


func (r *InMemoryOrderRepository) GetOrders(filter repository.OrderFilter) ([]models.Order, error) {
    // Implement the method to satisfy the OrderRepository interface
    return nil, nil
}

func (r *InMemoryOrderRepository) SaveOrder(order map[string]interface{}) error {
    // Implement the method to satisfy the OrderRepository interface
    return nil
}

func (r *InMemoryOrderRepository) UpdateOrderStatus(orderID, status string) error {
    // Implement the method to satisfy the OrderRepository interface
    return nil
}

func (r *InMemoryOrderRepository) GetOpenPositions() ([]models.Position, error) {
    // Implement the method to satisfy the OrderRepository interface
    return nil, nil
}

func (r *InMemoryOrderRepository) UpdatePosition(position models.Position) error {
    // Implement the method to satisfy the OrderRepository interface
    return nil
}

func (r *InMemoryOrderRepository) GetPosition(positionID string) (models.Position, error) {
    // Implement the method to satisfy the OrderRepository interface
    return models.Position{}, nil
}

func (r *InMemoryOrderRepository) ClosePosition(positionID string) error {
    // Implement the method to satisfy the OrderRepository interface
    return nil
}

// CreateScalperOrder processes high-frequency scalping orders with tight stop losses and quick profit-taking
func (s *OMSService) CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error) {
    if order.Price <= 0 || order.StopLoss <= 0 || order.RiskPercentage <= 0 {
        return nil, errors.New("invalid scalper order parameters")
    }

    order.ID = uuid.NewString()
    order.CreatedAt = time.Now().Unix()

    // Ensure quick execution and tight risk management
    if order.Price <= order.StopLoss {
        return nil, errors.New("price must be greater than stop loss")
    }

    // Calculate position size based on risk percentage
    accountBalance := 10000.0 // Example account balance
    riskAmount := accountBalance * order.RiskPercentage
    positionSize := riskAmount / (order.Price - order.StopLoss)
    order.Quantity = int(positionSize)

    return s.repo.CreateScalperOrder(order)
}

// ExecuteChildOrder manages the execution of child orders under a parent strategy (e.g., parent-child strategy)
func (s *OMSService) ExecuteChildOrder(parentID string) error {
    return s.repo.ExecuteChildOrder(parentID)
}

// GetTrades retrieves executed trades for a given parent order ID
func (s *OMSService) GetTrades(parentID string) ([]models.Trade, error) {
    return s.repo.GetTrades(parentID)
}

// CreateOrder creates a new order in the system (supports market, limit, and stop orders)
func (s *OMSService) CreateOrder(order models.Order) (*models.Order, error) {
    if order.Price <= 0 || order.Quantity <= 0 {
        return nil, errors.New("invalid order parameters")
    }

    order.ID = uuid.NewString()
    order.CreatedAt = time.Now().Unix()

    // Validate order type and strategy (e.g., Market, Limit, Stop)
    if order.Type == models.MarketOrder {
        // Direct market execution
        order.Status = models.OrderStatusExecuted
    } else if order.Type == models.LimitOrder {
        // Execute only if the market price matches the limit price
        order.Status = models.OrderStatusPending
    }

    createdOrder, err := s.repo.CreateOrder(order)
    if err != nil {
        return nil, err
    }
    return &createdOrder, nil
}
func (s *OMSService) GetOrders(filter repository.OrderFilter) ([]models.Order, error) {
    return s.repo.GetOrders(filter)
}

func (s *OMSService) ProcessOrder(order map[string]interface{}) error {
    if len(order) == 0 {
        return errors.New("invalid order data")
    }

    // Advanced validation logic for order parameters (e.g., price, quantity, stop loss)
    price, ok := order["price"].(float64)
    if !ok || price <= 0 {
        return errors.New("invalid order price")
    }

    quantity, ok := order["quantity"].(int)
    if !ok || quantity <= 0 {
        return errors.New("invalid order quantity")
    }

    orderType, ok := order["type"].(string)
    if !ok || (orderType != "market" && orderType != "limit" && orderType != "stop") {
        return errors.New("invalid order type")
    }

    // Additional validation for limit and stop orders
    if orderType == "limit" || orderType == "stop" {
        limitPrice, ok := order["limit_price"].(float64)
        if !ok || limitPrice <= 0 {
            return errors.New("invalid limit price")
        }
    }

    // Save the validated order to the database
    err := s.repo.SaveOrder(order)
    if err != nil {
        return err
    }

    // Execute the order if it's a market order
    if orderType == "market" {
        order["status"] = "executed"
        err = s.repo.UpdateOrderStatus(order["id"].(string), "executed")
        if err != nil {
            return err
        }
    } else {
        order["status"] = "pending"
    }

    return nil
}

// MonitorPositions periodically checks open positions and applies trailing stop losses and profit-taking
func (s *OMSService) MonitorPositions() error {
    positions, err := s.repo.GetOpenPositions()
    if err != nil {
        return err
    }

    for _, position := range positions {
        currentPrice := s.getCurrentPrice(position.Symbol) // Fetch real-time market price

        // Trailing stop logic
        if currentPrice > position.EntryPrice {
            newStopLoss := currentPrice - (position.EntryPrice - position.StopLoss)
            if newStopLoss > position.StopLoss {
                position.StopLoss = newStopLoss
            }
        }

        // Profit-taking strategy
        if currentPrice >= position.TakeProfit {
            if err := s.ClosePosition(position.ID); err != nil {
                return err
            }
        }

        position.CurrentPrice = currentPrice
        position.LastUpdatedAt = time.Now()
        s.repo.UpdatePosition(position)
    }

    return nil
}

// ClosePosition closes an open position by creating a sell order and updating the position status
func (s *OMSService) ClosePosition(positionID string) error {
    position, err := s.repo.GetPosition(positionID)
    if err != nil {
        return err
    }

    // Create a closing sell order for the position
    closingOrder := models.Order{
        ID:        uuid.NewString(),
        Symbol:    position.Symbol,
        Quantity:  position.Quantity,
        Price:     position.CurrentPrice,
        Side:      "sell",
        Type:      models.MarketOrder,
        Status:    models.OrderStatusPending,
        CreatedAt: time.Now().Unix(),
    }

    if _, err := s.CreateOrder(closingOrder); err != nil {
        return err
    }

    return s.repo.ClosePosition(positionID)
}

// getCurrentPrice simulates real-time market data retrieval
func (s *OMSService) getCurrentPrice(symbol string) float64 {
    // Replace this with actual market data API integration
    return 100.0
}

func ProcessOrder(order map[string]interface{}) error {
    if len(order) == 0 {
        return errors.New("invalid order data")
    }

    // Advanced validation logic for order parameters (e.g., price, quantity, stop loss)
    price, ok := order["price"].(float64)
    if !ok || price <= 0 {
        return errors.New("invalid order price")
    }

    quantity, ok := order["quantity"].(int)
    if !ok || quantity <= 0 {
        return errors.New("invalid order quantity")
    }

    orderType, ok := order["type"].(string)
    if !ok || (orderType != "market" && orderType != "limit" && orderType != "stop") {
        return errors.New("invalid order type")
    }

    // Additional validation for limit and stop orders
    if orderType == "limit" || orderType == "stop" {
        limitPrice, ok := order["limit_price"].(float64)
        if !ok || limitPrice <= 0 {
            return errors.New("invalid limit price")
        }
    }

    return nil
}
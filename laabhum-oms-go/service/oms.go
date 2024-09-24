package service

import (
	"errors"
	"time"
    "fmt"
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
    var filteredOrders []models.Order
    for _, order := range r.orders {
        if filter.Matches(order) {
            filteredOrders = append(filteredOrders, order)
        }
    }
    return filteredOrders, nil
}

func (r *InMemoryOrderRepository) SaveOrder(order map[string]interface{}) error {
    // Implement the method to satisfy the OrderRepository interface
    id, ok := order["id"].(string)
    if !ok {
        return errors.New("invalid order ID")
    }
    r.orders[id] = models.Order{
        ID:        id,
        Symbol:    order["symbol"].(string),
        Quantity:  order["quantity"].(int),
        Price:     order["price"].(float64),
        Type:      models.OrderType(order["type"].(string)),
        Status:    models.OrderStatus(order["status"].(string)),
        CreatedAt: order["created_at"].(int64),
    }
    return nil
}

func (r *InMemoryOrderRepository) UpdateOrderStatus(orderID, status string) error {
    // Implement the method to satisfy the OrderRepository interface
    order, ok := r.orders[orderID]
    if !ok {
        return errors.New("order not found")
    }
    order.Status = models.OrderStatus(status)
    r.orders[orderID] = order
    return nil
}

func (r *InMemoryOrderRepository) GetOpenPositions() ([]models.Position, error) {
    // Implement the method to satisfy the OrderRepository interface
    var openPositions []models.Position
    for _, position := range r.positions {
        if models.PositionStatus(position.Status) == models.PositionStatusOpen {
            openPositions = append(openPositions, position)
        }
    }
    return openPositions, nil
}

func (r *InMemoryOrderRepository) UpdatePosition(position models.Position) error {
    // Implement the method to satisfy the OrderRepository interface
    r.positions[position.ID] = position
    return nil
}

func (r *InMemoryOrderRepository) GetPosition(positionID string) (models.Position, error) {
    // Implement the method to satisfy the OrderRepository interface
    position, ok := r.positions[positionID]
    if !ok {
        return models.Position{}, errors.New("position not found")
    }
    return position, nil
}

func (r *InMemoryOrderRepository) ClosePosition(positionID string) error {
    // Implement the method to satisfy the OrderRepository interface
    position, ok := r.positions[positionID]
    if !ok {
        return errors.New("position not found")
    }
    position.Status = string(models.PositionStatusClosed)
    r.positions[positionID] = position
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
func (s *OMSService) ExecuteChildOrder(parentID, childID string) error {
    return s.repo.ExecuteChildOrder(childID)
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
    orders, err := s.repo.GetOrders(filter)
    if err != nil {
        return nil, err
    }
    var modelOrders []models.Order
    for _, order := range orders {
        modelOrders = append(modelOrders, models.Order{
            ID:        order.ID,
            Symbol:    order.Symbol,
            Quantity:  order.Quantity,
            Price:     order.Price,
            Type:      models.OrderType(order.Type),
            Status:    models.OrderStatus(order.Status),
            CreatedAt: order.CreatedAt,
        })
    }
    return modelOrders, nil
}

func (s *OMSService) validateOrderData(order map[string]interface{}) error {
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
    if (orderType == "limit" || orderType == "stop") {
        limitPrice, ok := order["limit_price"].(float64)
        if (!ok || limitPrice <= 0) {
            return errors.New("invalid limit price")
        }
    }

    return nil
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

func (s *OMSService) ExecuteAllChildTrades(parentID string) error {
    childOrders, err := s.repo.GetOrders(repository.OrderFilter{ParentID: parentID})
    if err != nil {
        return err
    }

    if len(childOrders) == 0 {
        return fmt.Errorf("no child orders found for parentID: %s", parentID) // Added check for no child orders
    }

    for _, childOrder := range childOrders {
        if err := s.ExecuteChildOrder(parentID, childOrder.ID); err != nil {
            return err // Handle execution error
        }
    }

    return nil
}


func (s *OMSService) ExecuteSpecificChild(parentID, childID string) error {
    // Implement the method to execute a specific child trade for a given parent order
    return s.repo.ExecuteChildOrder(childID)
}

func (s *OMSService) CreateCTC(order models.CTCOrder) (*models.CTCOrder, error) {
    // Implement the method to create a CTC order
    if order.Price <= 0 || order.Quantity <= 0 {
        return nil, errors.New("invalid CTC order parameters")
    }

    order.ID = uuid.NewString()
    order.CreatedAt = time.Now().Unix()

    _, err := s.repo.CreateOrder(models.Order{
        ID:        order.ID,
        Symbol:    order.Symbol,
        Quantity:  order.Quantity,
        Price:     order.Price,
        Type:      models.CTCOrderType,
        Status:    models.OrderStatusPending,
        CreatedAt: order.CreatedAt,
    })
    if err != nil {
        return nil, err
    }

    return &order, nil
}

func (s *OMSService) ExitAllTrades(parentID string) error {
    // Implement the method to exit all trades for a given parent order
    childOrders, err := s.repo.GetOrders(repository.OrderFilter{ParentID: parentID})
    if err != nil {
        return err
    }

    for _, childOrder := range childOrders {
        if err := s.repo.UpdateOrderStatus(childOrder.ID, models.OrderStatusCancelled); err != nil {
            return err
        }
    }

    return nil
}

func (s *OMSService) ExitChildTrades(parentID string) error {
    // Implement the method to exit all child trades for a given parent order
    return s.ExitAllTrades(parentID)
}

func (s *OMSService) ExitSpecificChild(parentID, childID string) error {
    // Implement the method to exit a specific child trade for a given parent order
    return s.repo.UpdateOrderStatus(childID, models.OrderStatusCancelled)
}

func (s *OMSService) CancelAllChildOrders(parentID string) error {
    // Implement the method to cancel all child orders for a given parent order
    childOrders, err := s.repo.GetOrders(repository.OrderFilter{ParentID: parentID})
    if err != nil {
        return err
    }

    for _, childOrder := range childOrders {
        if err := s.repo.UpdateOrderStatus(childOrder.ID, models.OrderStatusCancelled); err != nil {
            return err
        }
    }

    return nil
}

func (s *OMSService) CancelSpecificChildOrder(parentID, childID string) error {
    // Implement the method to cancel a specific child order for a given parent order
    return s.repo.UpdateOrderStatus(childID, models.OrderStatusCancelled)
}

func (s *OMSService) DeleteParentOrder(parentID string) error {
    // Implement the method to delete a parent order
    return s.repo.UpdateOrderStatus(parentID, models.OrderStatusDeleted)
}


func (s *OMSService) SyncPositions() error {
    // Implement the method to sync positions
    positions, err := s.repo.GetOpenPositions()
    if err != nil {
        return err
    }

    for _, position := range positions {
        currentPrice := s.getCurrentPrice(position.Symbol)
        position.CurrentPrice = currentPrice
        position.LastUpdatedAt = time.Now()
        if err := s.repo.UpdatePosition(position); err != nil {
            return err
        }
    }

    return nil
}

func (s *OMSService) ExecuteOrder(order models.Order) error {
    // Implement the method to execute an order
    if order.Type == models.MarketOrder {
        order.Status = models.OrderStatusExecuted
    } else {
        order.Status = models.OrderStatusPending
    }

    return s.repo.UpdateOrderStatus(order.ID, order.Status)
}

func (s *OMSService) CancelOrder(orderID string) error {
    // Implement the method to cancel an order
    return s.repo.UpdateOrderStatus(orderID, models.OrderStatusCancelled)
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
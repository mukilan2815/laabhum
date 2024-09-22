package models

import "time"

// Enums for Order Types, Status, and Trade Strategy
type OrderType string
type OrderStatus string
type TradeStrategy string

const (
    LimitOrder  OrderType = "LIMIT"
    MarketOrder OrderType = "MARKET"
    StopOrder   OrderType = "STOP"

    OrderStatusPending   OrderStatus = "PENDING"
    OrderStatusExecuted  OrderStatus = "EXECUTED"
    OrderStatusCancelled OrderStatus = "CANCELLED"

    StrategyDayTrading     TradeStrategy = "DAY_TRADING"
    StrategyPositionTrading TradeStrategy = "POSITION_TRADING"
    StrategyScalping        TradeStrategy = "SCALPING"
)

// Order represents a general order with advanced trading attributes
type Order struct {
    ID            string        `json:"id"`
    Symbol        string        `json:"symbol"`
    Quantity      int           `json:"quantity"`
    Price         float64       `json:"price"`
    Side          string        `json:"side"` // "buy" or "sell"
    Type          OrderType     `json:"type"` // LIMIT, MARKET, STOP
    Status        OrderStatus   `json:"status"` // PENDING, EXECUTED, CANCELLED
    StopPrice     float64       `json:"stop_price,omitempty"` // Stop Order Price (optional)
    Strategy      TradeStrategy `json:"strategy"` // Trading strategy (e.g. scalping, day trading)
    RiskPercentage float64      `json:"risk_percentage"` // % of capital risked
    StopLoss      float64       `json:"stop_loss"` // Stop loss level
    TakeProfit    float64       `json:"take_profit"` // Take profit level
    CreatedAt     int64         `json:"created_at"` // Timestamp for when the order is created
    ExpiresAt     time.Time     `json:"expires_at,omitempty"` // Optional expiry time for order
}

// Position represents an open position in the market
type Position struct {
    ID            string        `json:"id"`
    OrderID       string        `json:"order_id"`
    Symbol        string        `json:"symbol"`
    Quantity      int           `json:"quantity"`
    EntryPrice    float64       `json:"entry_price"` // Price at which the position was opened
    CurrentPrice  float64       `json:"current_price"` // Current market price
    StopLoss      float64       `json:"stop_loss"` // Dynamic stop-loss for trailing or fixed SL
    TakeProfit    float64       `json:"take_profit"` // Profit level to auto-close
    Strategy      TradeStrategy `json:"strategy"` // Associated trading strategy
    OpenedAt      time.Time     `json:"opened_at"` // Time when the position was opened
    LastUpdatedAt time.Time     `json:"last_updated_at"` // Last update timestamp for price/stop loss
}

// MarketCondition provides real-time or historical market data
type MarketCondition struct {
    Symbol     string    `json:"symbol"`
    Price      float64   `json:"price"` // Current price of the asset
    Volume     int       `json:"volume"` // Current market volume
    Volatility float64   `json:"volatility"` // Measure of price fluctuation
    Trend      string    `json:"trend"` // Market trend: bullish, bearish, sideways
    Timestamp  time.Time `json:"timestamp"` // Time of market data capture
}

// ScalperOrder represents a high-frequency order for scalping strategy
type ScalperOrder struct {
    ID           string    `json:"id"`
    Symbol       string    `json:"symbol"`
    Quantity     int       `json:"quantity"`
    Price        float64   `json:"price"` // Entry price for the scalper order
    StopLoss     float64   `json:"stop_loss"` // Tight stop-loss for scalping
    TakeProfit   float64   `json:"take_profit"` // Quick profit-taking level
    RiskPercentage float64 `json:"risk_percentage"` // % of capital at risk
    CreatedAt    int64     `json:"created_at"` // Timestamp for order creation
    ExpiresAt    time.Time `json:"expires_at,omitempty"` // Expiry time for the order
    Timestamp    string    `json:"timestamp"` // Timestamp for internal tracking
}

// Trade represents a successfully executed trade
type Trade struct {
    ID         string    `json:"id"`
    OrderID    string    `json:"order_id"` // The ID of the parent order that generated this trade
    Symbol     string    `json:"symbol"`
    Quantity   int       `json:"quantity"`
    Price      float64   `json:"price"` // Price at which the trade was executed
    TradeTime  time.Time `json:"trade_time"` // Time when the trade was executed
}

// Additional struct for handling historical data or advanced strategies
type HistoricalData struct {
    Symbol        string    `json:"symbol"`
    ClosePrices   []float64 `json:"close_prices"` // Array of closing prices for analysis
    Volatility    float64   `json:"volatility"`   // Historical volatility measure
    AverageVolume int       `json:"average_volume"` // Average volume over time
    Trend         string    `json:"trend"`        // Long-term trend based on price data
    Timestamp     time.Time `json:"timestamp"`    // Timestamp of the data point
}


package strategy

import (
	"errors"
	"log"
	"time"

	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
)

const (
	HighFrequencyTrader = "HighFrequencyTrader"
	Scalper             = "Scalper"
	DayTrader           = "DayTrader"
	PositionTrader      = "PositionTrader"
)

type Builder struct {
	logger      *log.Logger
	thresholds  map[string]int
	retryPolicy RetryPolicy
}

type RetryPolicy struct {
	MaxRetries int
	Delay      time.Duration
}

// NewBuilder creates a new Strategy Builder instance
func NewBuilder(logger *log.Logger, retryPolicy RetryPolicy) *Builder {
	return &Builder{
		logger:      logger,
		thresholds:  map[string]int{"high": 1000, "medium": 100, "low": 10},
		retryPolicy: retryPolicy,
	}
}

func (b *Builder) ProcessOrder(order oms.Order) (string, error) {
	b.logger.Printf("Processing order: %+v\n", order)

	// Implement strategy logic
	if order.Quantity <= 0 {
		err := errors.New("invalid order quantity")
		b.logger.Println(err)
		return "", err
	}

	// Determine the strategy based on quantity thresholds
	var strategy string
	if order.Quantity > b.thresholds["high"] {
		strategy = HighFrequencyTrader
	} else if order.Quantity > b.thresholds["medium"] {
		strategy = Scalper
	} else if order.Quantity > b.thresholds["low"] {
		strategy = DayTrader
	} else {
		strategy = PositionTrader
	}

	b.logger.Printf("Assigned strategy: %s\n", strategy)

	return strategy, nil
}

// RetryPolicy related functions
func (b *Builder) retryOperation(operation func() error) error {
	var lastError error
	for i := 0; i < b.retryPolicy.MaxRetries; i++ {
		lastError = operation()
		if lastError == nil {
			return nil
		}
		b.logger.Printf("Retrying operation, attempt %d/%d\n", i+1, b.retryPolicy.MaxRetries)
		time.Sleep(b.retryPolicy.Delay)
	}
	return lastError
}

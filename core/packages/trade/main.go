package trade

import (
	"fmt"
	"sort"
)

type ResponseStatus string

const (
	NoResponse ResponseStatus = "Open"
	Accepted   ResponseStatus = "Accepted"
	Declined   ResponseStatus = "Declined"
	Countered  ResponseStatus = "Countered"
)

type TradePlayerEntry struct {
	Status  ResponseStatus `json:"status"`
	Blocked bool           `json:"blocked"`
}

type TradeStatus string

const (
	TradeOpen      TradeStatus = "Open"
	TradeClosed    TradeStatus = "Closed"
	TradeFinalized TradeStatus = "Finalized"
)

type Trade struct {
	ID        int                          `json:"id"`
	Requester string                       `json:"requester"`
	Creator   string                       `json:"creator"`
	Responses map[string]*TradePlayerEntry `json:"responses"`
	Offer     map[string]int               `json:"offer"`
	Request   map[string]int               `json:"request"`
	Status    TradeStatus                  `json:"status"`
	ParentID  int                          `json:"parent"`
	Finalized bool                         `json:"finalized"`
	Timestamp int64                        `json:"timestamp"`
}

type Instance struct {
	// TODO: divide active trades with all trades
	// activeTrades     map[int]*Trade
	trades           map[int]*Trade
	parentToChildMap map[int][]int
	nextTradeID      int
}

func New() *Instance {
	return &Instance{
		// activeTrades:     make(map[int]*Trade),
		parentToChildMap: make(map[int][]int),
		nextTradeID:      1,
		trades:           make(map[int]*Trade),
	}
}

func (tm *Instance) Trades() []Trade {
	trades := make([]Trade, 0)
	for _, trade := range tm.trades {
		trades = append(trades, *trade)
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].ID < trades[j].ID
	})

	return trades
}

func (tm *Instance) ActiveTrades() []Trade {
	trades := make([]Trade, 0)
	for _, trade := range tm.trades {
		if trade.Status == TradeOpen {
			trades = append(trades, *trade)
		}
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].ID < trades[j].ID
	})

	return trades
}

func (tm *Instance) GetTrade(tradeID int) *Trade {
	return tm.trades[tradeID]
}

func (tm *Instance) CancelTrade(tradeID int) error {
	trade, exists := tm.trades[tradeID]
	if !exists {
		err := fmt.Errorf("Cannot cancel trade offer: trade#%d doesn't exist", tradeID)
		return err
	}

	// REFACTOR: probably redundant
	if trade.Finalized {
		err := fmt.Errorf("Cannot cancel trade offer: trade#%d already finalized", tradeID)
		return err
	}

	if trade.Status == TradeClosed {
		err := fmt.Errorf("Cannot cancel trade offer: trade#%d already closed", tradeID)
		return err
	}

	trade.Finalized = true
	trade.Status = TradeClosed
	return nil
}

func (tm *Instance) CancelActiveTrades() {
	for _, trade := range tm.trades {
		if trade.Status == TradeOpen {
			trade.Status = TradeClosed
		}
	}
}

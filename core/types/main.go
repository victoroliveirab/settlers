package types

import "fmt"

type HexCoordinate struct {
	Q int `json:"q"`
	R int `json:"r"`
	S int `json:"s"`
}

func (h HexCoordinate) String() string {
	return fmt.Sprintf("(q=%d,r=%d,s=%d)", h.Q, h.R, h.S)
}

// FIXME: redundant for now
type Settings struct {
	BankTradeAmount      int
	MaxCards             int
	MaxDevCardsPerRound  int
	MaxSettlements       int
	MaxCities            int
	MaxRoads             int
	TargetPoint          int
	PointsPerSettlement  int
	PointsPerCity        int
	PointsForMostKnights int
	PointsForLongestRoad int
	MostKnightsMinimum   int
	LongestRoadMinimum   int
}

type MapBlock struct {
	ID          int           `json:"id"`
	Resource    string        `json:"resource"`
	Token       int           `json:"token"`
	Vertices    [6]int        `json:"vertices"`
	Edges       [6]int        `json:"edges"`
	Coordinates HexCoordinate `json:"coordinates"`
	Blocked     bool          `json:"blocked"`
}

type Port struct {
	Type     string `json:"type"`
	Vertices [2]int `json:"vertices"`
}

type PlayerColor struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
}

type Player struct {
	ID    string      `json:"name"`
	Color PlayerColor `json:"color"`
}

type DevelopmentCard struct {
	Name        string `json:"name"`
	RoundBought int    `json:"-"`
}

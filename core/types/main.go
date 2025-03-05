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

type MapBlock struct {
	ID          int           `json:"id"`
	Resource    string        `json:"resource"`
	Token       int           `json:"token"`
	Vertices    [6]int        `json:"vertices"`
	Edges       [6]int        `json:"edges"`
	Coordinates HexCoordinate `json:"coordinates"`
	Blocked     bool          `json:"blocked"`
}

type Player struct {
	ID    string `json:"name"`
	Color string `json:"color"`
}

type DevelopmentCard struct {
	Name        string `json:"name"`
	RoundBought int    `json:"-"`
}

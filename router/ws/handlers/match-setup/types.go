package matchsetup

type newSettlementRequestPayload struct {
	VertexID int `json:"vertex"`
}

type newRoadRequestPayload struct {
	EdgeID int `json:"edge"`
}

package matchsetup

type currentRoundPlayerStateUpdateResponsePayload struct {
	Player string `json:"player"`
}

type verticesStateUpdateResponsePayload struct {
	AvailableVertices []int `json:"availableVertices"`
	Disabled          bool  `json:"disabled"`
}

type edgesStateUpdateResponsePayload struct {
	AvailableEdges []int `json:"availableEdges"`
	Disabled       bool  `json:"disabled"`
}

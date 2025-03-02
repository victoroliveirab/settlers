package matchsetup

import "fmt"

type settlementBuildPayload struct {
	vertexID int
}

type roadBuildPayload struct {
	edgeID int
}

func parseSettlementBuildPayload(payload map[string]interface{}) (*settlementBuildPayload, error) {
	vertexID, ok := payload["vertex"].(float64)
	if !ok {
		err := fmt.Errorf("malformed data: vertex")
		return nil, err
	}

	return &settlementBuildPayload{
		vertexID: int(vertexID),
	}, nil
}

func parseRoadBuildPayload(payload map[string]interface{}) (*roadBuildPayload, error) {
	edgeID, ok := payload["edge"].(float64)
	if !ok {
		err := fmt.Errorf("malformed data: edge")
		return nil, err
	}

	return &roadBuildPayload{
		edgeID: int(edgeID),
	}, nil
}

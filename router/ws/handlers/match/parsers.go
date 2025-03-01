package match

import "fmt"

type settlementBuildPayload struct {
	vertexID int
}

func parseSettlementBuildPayload(payload map[string]interface{}) (*settlementBuildPayload, error) {
	vertexID, ok := payload["vertex"].(int)
	if !ok {
		err := fmt.Errorf("malformed data: vertex")
		return nil, err
	}

	return &settlementBuildPayload{
		vertexID: vertexID,
	}, nil
}

package match

import "fmt"

type discardCardsPayload struct {
	resources map[string]int
}

func parseDiscardCardsPayload(payload map[string]interface{}) (*discardCardsPayload, error) {
	resourcesRaw, ok := payload["resources"]
	if !ok {
		return nil, fmt.Errorf("missing 'resources' key in payload")
	}

	resources := make(map[string]int)

	for key, value := range resourcesRaw.(map[string]interface{}) {
		num, ok := value.(float64)
		if !ok {
			err := fmt.Errorf("malformed data for key %s", key)
			return nil, err
		}
		resources[key] = int(num)
	}

	return &discardCardsPayload{resources: resources}, nil
}

type roadBuildPayload struct {
	edgeID int
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

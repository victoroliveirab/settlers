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

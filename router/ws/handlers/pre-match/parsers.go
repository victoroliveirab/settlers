package prematch

import "fmt"

type updateParamPayload struct {
	Key   string
	Value int
}

func parseUpdateParamPayload(payload map[string]interface{}) (*updateParamPayload, error) {
	key, ok := payload["key"].(string)
	if !ok {
		err := fmt.Errorf("malformed data: key")
		return nil, err
	}

	value, ok := payload["value"].(float64)
	if !ok {
		err := fmt.Errorf("malformed data: value")
		return nil, err
	}
	return &updateParamPayload{
		Key:   key,
		Value: int(value),
	}, nil
}

type roomPayload struct {
	RoomID string
}

func parseRoomJoinPayload(payload map[string]interface{}) (*roomPayload, error) {
	roomID, ok := payload["roomID"].(string)
	if !ok {
		err := fmt.Errorf("malformed data: roomID")
		return nil, err
	}

	return &roomPayload{
		RoomID: roomID,
	}, nil
}

type playerReadyPayload struct {
	Ready  bool
	RoomID string
}

func parsePlayerReadyState(payload map[string]interface{}) (*playerReadyPayload, error) {
	ready, ok := payload["ready"].(bool)
	if !ok {
		err := fmt.Errorf("malformed data: ready")
		return nil, err
	}
	roomID, ok := payload["roomID"].(string)
	if !ok {
		err := fmt.Errorf("malformed data: roomID")
		return nil, err
	}

	return &playerReadyPayload{
		Ready:  ready,
		RoomID: roomID,
	}, nil
}

package prematch

import "fmt"

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

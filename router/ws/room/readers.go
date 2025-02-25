package room

import "fmt"

type roomPayload struct {
	RoomID string
}

func readRoomPayload(payload map[string]interface{}) (*roomPayload, error) {
	roomID, ok := payload["roomID"].(string)
	if !ok {
		err := fmt.Errorf("malformed data: roomID")
		return nil, err
	}

	return &roomPayload{
		RoomID: roomID,
	}, nil
}

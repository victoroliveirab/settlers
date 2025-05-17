package entities

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/core/packages/round"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/utils"
)

var availableColors []coreT.PlayerColor = []coreT.PlayerColor{
	coreT.PlayerColor{
		Background: "palegreen",
		Foreground: "black",
	},
	coreT.PlayerColor{
		Background: "orange",
		Foreground: "black",
	},
	coreT.PlayerColor{
		Background: "maroon",
		Foreground: "white",
	},
	coreT.PlayerColor{
		Background: "lemonchiffon",
		Foreground: "black",
	},
	coreT.PlayerColor{
		Background: "blue",
		Foreground: "white",
	},
	coreT.PlayerColor{
		Background: "crimson",
		Foreground: "white",
	},
	coreT.PlayerColor{
		Background: "orangered",
		Foreground: "black",
	},
	coreT.PlayerColor{
		Background: "aliceblue",
		Foreground: "black",
	},
	coreT.PlayerColor{
		Background: "lightslategray",
		Foreground: "black",
	},
}

func NewRoom(id, mapName string, capacity, randSeed int, params RoomParams, onDestroy func(room *Room)) *Room {
	randGenerator := utils.RandNew(int64(randSeed))
	ctx, cancel := context.WithCancel(context.Background())
	return &Room{
		ID:               id,
		Capacity:         capacity,
		MapName:          mapName,
		Participants:     make([]RoomEntry, capacity),
		Owner:            "",
		Colors:           availableColors,
		params:           params,
		Status:           "prematch",
		incomingMsgQueue: make(chan IncomingMessage, 32), // buffer incoming messages
		outgoingMsgQueue: make(chan OutgoingMessage),     // process msg immediatly, one by one
		handlers:         make([]RoomIncomingMessageHandler, 0),
		onDestroy:        onDestroy,
		Rand:             randGenerator,
		Game:             nil,
		MaxIdleTime:      5 * time.Second,
		Private:          true,
		ctx:              ctx,
		cancel:           cancel,
	}
}

func (room *Room) CanAddPlayer() bool {
	if room.Status != "prematch" {
		fmt.Println("CANNOT JOIN ROOM BECAUSE ROOM IS", room.Status)
		return false
	}
	for _, participant := range room.Participants {
		if participant.Player == nil {
			return true
		}
	}
	return false
}

func (room *Room) AddPlayer(player *GamePlayer) error {
	room.Lock()
	defer room.Unlock()

	if room.destroyTimerCallback != nil {
		logger.LogSystemMessage("room.CancelDestroy", fmt.Sprintf("Cancelling scheduled destruction for room %s due to player %s joining", room.ID, player.Username))
		room.destroyTimerCallback.Stop()
		room.destroyTimerCallback = nil
	}

	for _, spot := range room.Participants {
		if spot.Player != nil {
			if spot.Player.ID == player.ID {
				err := fmt.Errorf("Cannot add player %s to room#%s: already in room", player.Username, room.ID)
				return err
			}
		}
	}
	for i, spot := range room.Participants {
		if spot.Player == nil {
			player.Color = &availableColors[i]
			room.Participants[i] = RoomEntry{
				Player: player,
				Ready:  false,
				Bot:    false,
			}
			if room.Owner == "" {
				room.Owner = player.Username
			}
			return nil
		}
	}

	err := fmt.Errorf("Cannot join room #%s: room full", room.ID)
	return err
}

func (room *Room) ChangePlayerColor(playerID int64, color string) error {
	room.Lock()
	defer room.Unlock()

	if colorIndex := utils.SliceFindIndex(availableColors, func(val coreT.PlayerColor) bool { return val.Background == color }); colorIndex == -1 {
		err := fmt.Errorf("Cannot use color %s: unknown color", color)
		return err
	}

	for _, participant := range room.Participants {
		if participant.Player != nil && participant.Player.Color.Background == color {
			err := fmt.Errorf("Cannot use color %s: color taken", color)
			return err
		}
	}
	for index, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == playerID {
			room.Participants[index].Player.Color = &availableColors[index]
			return nil
		}
	}

	err := fmt.Errorf("Cannot set color for %d: player not found", playerID)
	return err
}

func (room *Room) ReconnectPlayer(player *GamePlayer) error {
	room.Lock()
	defer room.Unlock()
	for i := range room.Participants {
		if room.Participants[i].Player != nil && room.Participants[i].Player.ID == player.ID {
			room.Participants[i].Bot = false
			break
		}
	}

	err := fmt.Errorf("Cannot remove player#%d: not part of room %s", player.ID, room.ID)
	return err

}

// THIS FUNCTION HAS TOO MANY RESPONSIBILITIES
// MAKE THIS JUST REMOVE PLAYER FROM ROOM
// LET ASSIGN NEW OWNER TO SOMEWHERE ELSE
// AND ALSO THE SCHEDULING OF DESTRUCTION
func (room *Room) RemovePlayer(playerID int64) error {
	fmt.Println("REMOVE PLAYER:", playerID)
	debug.PrintStack()
	room.Lock()
	defer room.Unlock()

	var participantIndex int = -1
	for i := range room.Participants {
		if room.Participants[i].Player != nil && room.Participants[i].Player.ID == playerID {
			participantIndex = i
			break
		}
	}

	if participantIndex == -1 {
		err := fmt.Errorf("Cannot remove player#%d: not part of room %s", playerID, room.ID)
		fmt.Println(err.Error())
		return err
	}

	participant := room.Participants[participantIndex]
	//participant.Player.Room = nil
	participantName := participant.Player.Username

	shouldDestroyRoom := false
	var messageType string
	var log string

	if room.Status == "prematch" || room.Game == nil {
		room.Participants[participantIndex] = RoomEntry{}
		messageType = "room.player-left"
		log = fmt.Sprintf("%s has left the room", participantName)
	} else {
		room.Participants[participantIndex].Player.Disconnect()
		room.Participants[participantIndex].Bot = true
		messageType = "match.player-left"
		log = fmt.Sprintf("%s has left the match. A bot was assigned in their place.", participantName)
	}

	isOwner := participantName == room.Owner
	if isOwner {
		err := room.assignNewOwner()
		if err != nil {
			shouldDestroyRoom = true
		}
	}

	fmt.Println("WILL SCHEDULE DESTROY ROOM", shouldDestroyRoom)

	if shouldDestroyRoom {
		room.scheduleRoomDestroy()
	} else {
		recipients := []string{}
		for _, participant := range room.Participants {
			if participant.Player != nil && !participant.Bot && participant.Player.GetConnection() != nil {
				recipients = append(recipients, participant.Player.Username)
			}
		}
		room.EnqueueOutgoingMessage(&types.WebSocketServerResponse{
			Type: types.ResponseType(messageType),
			Payload: map[string]interface{}{
				"player": participantName,
				"logs":   []string{log},
			},
		}, recipients, nil)
	}
	return nil
}

func (room *Room) scheduleRoomDestroy() {
	if room.destroyTimerCallback != nil {
		room.destroyTimerCallback.Stop()
	}
	logger.LogSystemMessage("room.scheduleRoomDestroy", fmt.Sprintf("Scheduling destruction for room %s in 15s due to last player leaving", room.ID))
	cb := time.AfterFunc(5*time.Second, func() {
		room.destroy(fmt.Sprintf("Room %s destroyed after 15s of inactivity (owner left)", room.ID))
	})
	room.destroyTimerCallback = cb
}

func (room *Room) ProgressStatus() error {
	if room.Status == "prematch" {
		room.Status = "setup"
		return nil
	}
	if room.Status == "setup" {
		room.Status = "match"
		room.StartDatetime = time.Now().UTC()
		return nil
	}
	if room.Status == "match" {
		room.Status = "over"
		room.EndDatetime = time.Now().UTC()
		return nil
	}
	err := fmt.Errorf("Cannot proceed status %s", room.Status)
	return err
}

func (room *Room) assignNewOwner() error {
	for _, participant := range room.Participants {
		if participant.Player != nil && participant.Player.GetConnection() != nil && !participant.Bot {
			room.Owner = participant.Player.Username
			fmt.Println("NEW OWNER:", room.Owner)
			return nil
		}
	}
	err := fmt.Errorf("Cannot assign a new owner to room %s: no players left", room.ID)
	return err
}

func (room *Room) CreateRoundManager(onTimeout func(), onExpireFuncs map[round.Type]func()) error {
	if room.roundManager != nil {
		err := fmt.Errorf("Error: room#%s already has a round manager initialized", room.ID)
		return err
	}

	room.roundManager = newRoundManager(room.params.Values["speed"], onTimeout, onExpireFuncs)
	return nil
}

func (room *Room) StartRound() error {
	if room.roundManager == nil {
		err := fmt.Errorf("Error: room#%s doesn't have a round manager initialized", room.ID)
		return err
	}
	room.roundManager.start()
	return nil
}

func (room *Room) ResumeRound() error {
	if room.roundManager == nil {
		err := fmt.Errorf("Error: room#%s doesn't have a round manager initialized", room.ID)
		return err
	}
	room.roundManager.cancelSubTimer()
	room.roundManager.resume()
	return nil
}

func (room *Room) StartSubRound(phase round.Type) error {
	if room.roundManager == nil {
		err := fmt.Errorf("Error: room#%s doesn't have a round manager initialized", room.ID)
		return err
	}
	room.roundManager.pause()
	room.roundManager.startPhaseTimer(phase)
	return nil
}

func (room *Room) EndRound() error {
	if room.roundManager == nil {
		err := fmt.Errorf("Error: room#%s doesn't have a round manager initialized", room.ID)
		return err
	}
	room.roundManager.cancel()
	return nil

}

func (room *Room) Now() time.Time {
	return room.roundManager.Now()
}

func (room *Room) RoundDeadline() *time.Time {
	return room.roundManager.Deadline()
}

func (room *Room) SubRoundDeadline() *time.Time {
	return room.roundManager.SubPhaseDeadline()
}

func (room *Room) RegisterIncomingMessageHandler(f RoomIncomingMessageHandler) {
	room.handlers = append(room.handlers, f)
}

func (room *Room) EnqueueIncomingMessage(player *GamePlayer, msg *types.WebSocketClientRequest) {
	room.incomingMsgQueue <- IncomingMessage{
		Player:  player,
		Message: msg,
	}
}

func (room *Room) EnqueueBulkUpdate(updaters ...func(room *Room, username string) *types.WebSocketServerResponse) {
	messages := make([]OutgoingMessage, room.Capacity)
	messageType := fmt.Sprintf("%s.bulk-update", room.Status)
	for i, participant := range room.Participants {
		if participant.Player != nil {
			messages[i] = OutgoingMessage{
				Message: &types.WebSocketServerResponse{
					Type:    types.ResponseType(messageType),
					Payload: []types.WebSocketServerResponse{},
				},
				Recipients: []string{participant.Player.Username},
			}
		}
	}

	for _, updater := range updaters {
		for index, participant := range room.Participants {
			if participant.Player != nil {
				update := updater(room, participant.Player.Username)
				if msg, ok := messages[index].Message.Payload.([]types.WebSocketServerResponse); ok {
					messages[index].Message.Payload = append(msg, *update)
				}
			}
		}
	}

	for _, message := range messages {
		room.outgoingMsgQueue <- message
	}
}

func (room *Room) EnqueueOutgoingMessage(msg *types.WebSocketServerResponse, recipients []string, onSend func()) {
	room.outgoingMsgQueue <- OutgoingMessage{
		Callback:   onSend,
		Message:    msg,
		Recipients: recipients,
	}
}

func (room *Room) ProcessIncomingMessages() {
	for {
		select {
		case <-room.ctx.Done():
			fmt.Println("DONE PROCESSING INCOMING MESSAGES")
			return
		case item := <-room.incomingMsgQueue:
			message := item.Message
			sender := item.Player

			sender.SetLastActiveTimestamp(time.Now())

			var handled bool
			var err error

			for _, handler := range room.handlers {
				handled, err = handler(sender, message)
				if handled || err != nil {
					break
				}
			}
			if handled && err == nil {
				continue
			}

			if !handled {
				err = fmt.Errorf("Unknown message type: %s", message.Type)
			}
			logger.LogError(sender.ID, fmt.Sprintf("room-%s.ProcessIncomingMessages", room.ID), -1, err)

		}
	}
}

func (room *Room) ProcessOutgoingMessages() {
	for {
		select {
		case <-room.ctx.Done():
			return
		case item := <-room.outgoingMsgQueue:
			var recipients []RoomEntry

			// Copy so if a participant disconnects mid broadcast and the room.Participants array changes, we don't panic
			room.Lock()
			for _, participant := range room.Participants {
				if len(item.Recipients) == 0 {
					recipients = append(recipients, participant)
					continue
				}
				for _, recipient := range item.Recipients {
					if participant.Player != nil && participant.Player.Username == recipient {
						recipients = append(recipients, participant)
					}
				}
			}
			room.Unlock()

			for _, participant := range recipients {
				player := participant.Player
				if player == nil || player.GetConnection() == nil {
					continue
				}
				wsErr := player.WriteJSON(item.Message)
				if wsErr != nil {
					logger.LogError(player.ID, fmt.Sprintf("room-%s.ProcessOutgoingMessages", room.ID), -1, wsErr)
					if !websocket.IsCloseError(wsErr) && !websocket.IsUnexpectedCloseError(wsErr) {
						logger.LogSystemMessage(
							fmt.Sprintf("room-%s.ProcessOutgoingMessages.%d.OnDisconnect", room.ID, player.ID),
							fmt.Sprintf("rrror wasn't due to trying to write in closed connection"),
						)
						go player.OnDisconnect(player)
					}
					continue
				}
			}

			if item.Callback != nil {
				go item.Callback()
			}
		}
	}
}

func (room *Room) destroy(reason string) {
	logger.LogSystemMessage("room.Destroy", reason)
	room.roundManager.cancel()

	room.roundManager = nil // Mark roundManager as garbage collectable
	room.Game = nil         // Mark Game as garbage collectable
	room.cancel()

	for _, participant := range room.Participants {
		if participant.Player != nil {
			participant.Player.Destroy()
			participant.Player = nil
		}
	}

	room.onDestroy(room)
}

func (room *Room) MinMax() [2]int {
	var minMax [2]int
	minMax[0] = room.params.MinPlayers
	minMax[1] = room.params.MaxPlayers
	return minMax
}

func (room *Room) Params() []RoomParamsMetaEntry {
	var entries []RoomParamsMetaEntry
	for _, v := range room.params.Meta {
		entries = append(entries, RoomParamsMetaEntry{
			Key:         v.Key,
			Description: v.Description,
			Label:       v.Label,
			Priority:    v.Priority,
			Value:       room.params.Values[v.Key],
			Values:      v.Values,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Priority == entries[j].Priority {
			return entries[i].Key < entries[j].Key // Alphabetical order
		}
		return entries[i].Priority > entries[j].Priority // Higher priority first
	})

	return entries
}

func (room *Room) UpdateSize(player *GamePlayer, newSize int) error {
	if room.Owner != player.Username {
		err := fmt.Errorf("cannot update size in room %s: not room owner", room.ID)
		return err
	}

	room.Lock()
	defer room.Unlock()

	currentNumberOfParticipants := 0
	for _, spot := range room.Participants {
		if spot.Player != nil {
			currentNumberOfParticipants++
		}
	}
	if newSize < currentNumberOfParticipants {
		err := fmt.Errorf("cannot shrink room size to %d in room %s: too many players", newSize, room.ID)
		return err
	}
	participants := make([]RoomEntry, newSize)
	currentIndex := 0
	for _, spot := range room.Participants {
		if spot.Player != nil {
			participants[currentIndex] = spot
			currentIndex++
		}
	}
	room.Participants = participants
	room.Capacity = newSize
	return nil
}

func (room *Room) UpdateParam(player *GamePlayer, key string, value int) error {
	if room.Owner != player.Username {
		err := fmt.Errorf("cannot update param %s in room %s: not room owner", key, room.ID)
		return err
	}

	_, exists := room.params.Values[key]
	if !exists {
		err := fmt.Errorf("unknown param: %s", key)
		return err
	}
	if !utils.SliceContains(room.params.Meta[key].Values, value) {
		err := fmt.Errorf("invalid value for param %s: %d", key, value)
		return err
	}
	room.params.Values[key] = value
	return nil
}

func (room *Room) TogglePlayerReadyState(playerID int64, newState bool) error {
	room.Lock()
	defer room.Unlock()

	for index, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == playerID {
			room.Participants[index].Ready = newState
			return nil
		}
	}

	err := fmt.Errorf("Cannot toggle player#%d ready state to %v: not part of room %s", playerID, newState, room.ID)
	return err
}

package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type makeBankTradeRequestPayload struct {
	Given     map[string]int `json:"given"`
	Requested map[string]int `json:"requested"`
}

type makeGeneralPortTradeRequestPayload struct {
	Given     map[string]int `json:"given"`
	Requested map[string]int `json:"requested"`
}

type makeResourcePortTradeRequestPayload struct {
	Given     map[string]int `json:"given"`
	Requested map[string]int `json:"requested"`
}

type createTradeOfferRequestPayload struct {
	Given     map[string]int `json:"given"`
	Requested map[string]int `json:"requested"`
}

type createCounterTradeOfferRequestPayload struct {
	Given     map[string]int `json:"given"`
	Requested map[string]int `json:"requested"`
	TradeID   int            `json:"tradeID"`
}

type acceptTradeOfferRequestPayload struct {
	TradeID int `json:"tradeID"`
}

type rejectTradeOfferRequestPayload struct {
	TradeID int `json:"tradeID"`
}

type finalizeTradeOfferRequestPayload struct {
	AccepterID string `json:"accepter"`
	TradeID    int    `json:"tradeID"`
}

type cancelTradeOfferRequestPayload struct {
	TradeID int `json:"tradeID"`
}

func handleMakeBankTrade(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[makeBankTradeRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	resourceGiven := payload.Given
	resourceRequested := payload.Requested
	room := player.Room
	game := room.Game

	err = game.MakeBankTrade(player.Username, resourceGiven, resourceRequested)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	formattedResourceGiven := formatResourceCollection(resourceGiven)
	formattedResourceRequested := formatResourceCollection(resourceRequested)

	logs := []string{fmt.Sprintf("%s traded %s for %s with the bank", player.Username, formattedResourceGiven, formattedResourceRequested)}
	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateLogs(logs),
	)

	return true, nil
}

func handleMakeGeneralPortTrade(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[makeGeneralPortTradeRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	resourceGiven := payload.Given
	resourceRequested := payload.Requested
	room := player.Room
	game := room.Game

	err = game.MakeGeneralPortTrade(player.Username, resourceGiven, resourceRequested)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}
	formattedResourceGiven := formatResourceCollection(resourceGiven)
	formattedResourceRequested := formatResourceCollection(resourceRequested)

	logs := []string{fmt.Sprintf("%s traded %s for %s with the port", player.Username, formattedResourceGiven, formattedResourceRequested)}
	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateLogs(logs),
	)

	return true, nil
}

func handleMakeResourcePortTrade(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[makeResourcePortTradeRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	resourceGiven := payload.Given
	resourceRequested := payload.Requested
	room := player.Room
	game := room.Game

	err = game.MakeResourcePortTrade(player.Username, resourceGiven, resourceRequested)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	formattedResourceGiven := formatResourceCollection(resourceGiven)
	formattedResourceRequested := formatResourceCollection(resourceRequested)

	logs := []string{fmt.Sprintf("%s traded %s for %s with the port", player.Username, formattedResourceGiven, formattedResourceRequested)}
	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateLogs(logs),
	)

	return true, nil

}

func handleCreateTradeOffer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[createTradeOfferRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	resourcesGiven := payload.Given
	resourcesRequested := payload.Requested
	room := player.Room
	game := room.Game

	_, err = game.MakeTradeOffer(player.Username, resourcesGiven, resourcesRequested, []string{})
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	logs := make([]string, 1)
	logs[0] = fmt.Sprintf("%s is offering to trade", player.Username)

	room.EnqueueBulkUpdate(
		UpdateTradeOffers,
		UpdateLogs(logs),
	)

	return true, nil
}

func handleCreateCounterTradeOffer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[createCounterTradeOfferRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	tradeID := payload.TradeID
	resourcesGiven := payload.Given
	resourcesRequested := payload.Requested
	room := player.Room
	game := room.Game

	_, err = game.MakeCounterTradeOffer(player.Username, tradeID, resourcesGiven, resourcesRequested)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room.EnqueueBulkUpdate(
		UpdateTradeOffers,
	)

	return true, nil
}

func handleAcceptTradeOffer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[acceptTradeOfferRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	tradeID := payload.TradeID
	room := player.Room
	game := room.Game
	err = game.AcceptTradeOffer(player.Username, tradeID)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room.EnqueueBulkUpdate(
		UpdateTradeOffers,
	)

	return true, nil
}

func handleRejectTradeOffer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[rejectTradeOfferRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	tradeID := payload.TradeID
	room := player.Room
	game := room.Game
	err = game.RejectTradeOffer(player.Username, tradeID)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room.EnqueueBulkUpdate(
		UpdateTradeOffers,
	)

	return true, nil
}

func handleFinalizeTradeOffer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[finalizeTradeOfferRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	accepterID := payload.AccepterID
	tradeID := payload.TradeID
	room := player.Room
	game := room.Game
	err = game.FinalizeTrade(player.Username, accepterID, tradeID)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	trade := game.GetTradeByID(tradeID)

	formattedResourceGiven := formatResourceCollection(trade.Offer)
	formattedResourceRequested := formatResourceCollection(trade.Request)

	logs := []string{fmt.Sprintf("%s traded %s for %s with %s", player.Username, formattedResourceGiven, formattedResourceRequested, accepterID)}
	room.EnqueueBulkUpdate(
		UpdatePlayerHand,
		UpdateResourceCount,
		UpdateTradeOffers,
		UpdateBuyDevelopmentCard,
		UpdateLogs(logs),
	)

	return true, nil
}

func handleCancelTradeOffer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[cancelTradeOfferRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	tradeID := payload.TradeID
	room := player.Room
	game := room.Game
	err = game.CancelTradeOffer(player.Username, tradeID)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room.EnqueueBulkUpdate(
		UpdateTradeOffers,
	)

	return true, nil
}

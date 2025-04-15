package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	switch message.Type {
	case "room.update-capacity":
		requestPayload, err := utils.ParseJsonPayload[roomUpdateCapacityRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		err = room.UpdateSize(player, requestPayload.Capacity)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.update-param":
		requestPayload, err := utils.ParseJsonPayload[roomUpdateParamRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		key := requestPayload.Key
		value := requestPayload.Value
		err = room.UpdateParam(player, key, value)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.player-change-color":
		requestPayload, err := utils.ParseJsonPayload[roomPlayerChangeColorRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		color := requestPayload.Color
		err = room.ChangePlayerColor(player.ID, color)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.toggle-ready":
		requestPayload, err := utils.ParseJsonPayload[roomPlayerReadyRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		ready := requestPayload.Ready
		err = room.TogglePlayerReadyState(player.ID, ready)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.start-game":
		room := player.Room
		err := StartMatch(player, room)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		// game := room.Game
		// err = game.BuildSettlement("victoroliveirab", 2)
		// game.BuildRoad("victoroliveirab", 2)
		//
		// game.BuildSettlement("teste", 28)
		// game.BuildRoad("teste", 37)
		//
		// game.BuildSettlement("teste", 33)
		// game.BuildRoad("teste", 43)
		//
		// game.BuildSettlement("victoroliveirab", 8)
		// game.BuildRoad("victoroliveirab", 9)

		// game.BuildSettlement("victoroliveirab", 2)
		// game.BuildRoad("victoroliveirab", 2)
		//
		// game.BuildSettlement("coxa", 28)
		// game.BuildRoad("coxa", 37)
		//
		// game.BuildSettlement("teste", 45)
		// game.BuildRoad("teste", 60)
		//
		// game.BuildSettlement("barbosa", 22)
		// game.BuildRoad("barbosa", 28)
		//
		// game.BuildSettlement("barbosa", 24)
		// game.BuildRoad("barbosa", 29)
		//
		// game.BuildSettlement("teste", 8)
		// game.BuildRoad("teste", 9)
		//
		// game.BuildSettlement("coxa", 33)
		// game.BuildRoad("coxa", 43)
		//
		// game.BuildSettlement("victoroliveirab", 4)
		// game.BuildRoad("victoroliveirab", 3)
		//
		// room.ProgressStatus()
		//
		// err = game.RollDice("victoroliveirab")
		// game.EndRound("victoroliveirab")
		//
		// game.RollDice("coxa")
		// game.BuildRoad("coxa", 44)
		// game.EndRound("coxa")
		//
		// game.RollDice("teste")
		// game.EndRound("teste")
		//
		// game.RollDice("barbosa")
		// tradeID, _ := game.MakeTradeOffer("barbosa", map[string]int{
		// 	"Grain": 1,
		// 	"Ore":   1,
		// }, map[string]int{
		// 	"Sheep": 1,
		// }, []string{})
		// err = game.AcceptTradeOffer("victoroliveirab", tradeID)
		// if err != nil {
		// 	panic(err)
		// }
		// err = game.FinalizeTrade("barbosa", "victoroliveirab", tradeID)
		// if err != nil {
		// 	panic(err)
		// }
		// err = game.BuyDevelopmentCard("barbosa")
		// if err != nil {
		// 	panic(err)
		// }
		// game.EndRound("barbosa")
		//
		// game.RollDice("victoroliveirab")
		// game.EndRound("victoroliveirab")
		//
		// game.RollDice("coxa")
		// tradeID, _ = game.MakeTradeOffer("coxa", map[string]int{
		// 	"Brick": 1,
		// }, map[string]int{
		// 	"Ore": 1,
		// }, []string{})
		// game.AcceptTradeOffer("victoroliveirab", tradeID)
		// game.FinalizeTrade("coxa", "victoroliveirab", tradeID)
		// game.BuildCity("coxa", 28)
		// game.EndRound("coxa")
		//
		// game.RollDice("teste")
		// game.BuildRoad("teste", 72)
		// tradeID, _ = game.MakeTradeOffer("teste", map[string]int{
		// 	"Lumber": 1,
		// 	"Grain":  1,
		// }, map[string]int{
		// 	"Brick": 1,
		// }, []string{})
		// err = game.AcceptTradeOffer("victoroliveirab", tradeID)
		// if err != nil {
		// 	panic(err)
		// }
		// err = game.FinalizeTrade("teste", "victoroliveirab", tradeID)
		// if err != nil {
		// 	panic(err)
		// }
		// err = game.BuildSettlement("teste", 54)
		// if err != nil {
		// 	panic(err)
		// }
		// game.EndRound("teste")
		//
		// game.UseKnight("barbosa")
		// game.MoveRobber("barbosa", 6)
		// game.RobPlayer("barbosa", "teste")
		// game.RollDice("barbosa")
		// game.MoveRobber("barbosa", 1)
		// game.RobPlayer("barbosa", "victoroliveirab")
		// game.EndRound("barbosa")
		//
		// game.RollDice("victoroliveirab")
		// game.BuildCity("victoroliveirab", 4)
		// game.EndRound("victoroliveirab")
		//
		// game.RollDice("coxa")
		// game.EndRound("coxa")
		//
		// game.RollDice("teste")
		// game.EndRound("teste")
		//
		// game.RollDice("barbosa")
		// game.EndRound("barbosa")
		//
		// game.RollDice("victoroliveirab")
		// game.MoveRobber("victoroliveirab", 15)
		// game.RobPlayer("victoroliveirab", "teste")
		// game.EndRound("victoroliveirab")
		//
		// game.RollDice("coxa")
		// game.MakeBankTrade("coxa", "Grain", "Sheep")
		// game.BuyDevelopmentCard("coxa")
		// game.EndRound("coxa")
		//
		// game.RollDice("teste")
		// game.EndRound("teste")
		//
		// game.RollDice("barbosa")
		// game.BuildCity("barbosa", 22)
		// game.EndRound("barbosa")
		//
		// game.RollDice("victoroliveirab")
		// game.EndRound("victoroliveirab")

		room.EnqueueOutgoingMessage(buildStartMatch(room), nil, func() {
			room.EnqueueBulkUpdate(
				match.UpdateCurrentRoundPlayerState,
				match.UpdateVertexState,
				match.UpdateEdgeState,
				match.UpdateLogs([]string{"Setup phase starting."}),
			)
		})
		return true, nil
	default:
		return false, nil
	}
}

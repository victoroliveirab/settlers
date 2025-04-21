package core

// import (
// 	"testing"
// )

// func TestBuildCitySuccess(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithSettlementsByPlayer(map[string][]int{
// 			"1": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  3,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build success", func(t *testing.T) {
// 		settlements := game.AllSettlements()
// 		_, exists := settlements[42]
// 		if !exists {
// 			t.Errorf("expected settlement to exist at vertex#42 before city built, but it actually doesn't")
// 		}
//
// 		err := game.BuildCity("1", 42)
// 		if err != nil {
// 			t.Errorf("expected to be able to build city in vertex#42 during regular phase, but found error %s", err.Error())
// 		}
//
// 		settlements = game.AllSettlements()
// 		settlement := settlements[42]
// 		emptyBuilding := Building{
// 			Owner: "",
// 			ID:    0,
// 		}
// 		if settlement != emptyBuilding {
// 			t.Errorf("expected settlement to not exist at vertex#42 after city built, but it actually does")
// 		}
//
// 		cities := game.AllCities()
// 		newCity := cities[42]
// 		if newCity == emptyBuilding {
// 			t.Errorf("expected new city to show up in cities map, but it didn't")
// 		}
// 		if newCity.Owner != "1" {
// 			t.Errorf("expected new city to belong to player#1, but it actually belongs to %s", newCity.Owner)
// 		}
//
// 		player1ResourcesAfterBuild := game.ResourceHandByPlayer("1")
// 		if player1ResourcesAfterBuild["Lumber"] != 0 {
// 			t.Errorf("expected to have 0 Lumber after build city, but found %d", player1ResourcesAfterBuild["Lumber"])
// 		}
// 		if player1ResourcesAfterBuild["Brick"] != 1 {
// 			t.Errorf("expected to have 2 Brick after build city, but found %d", player1ResourcesAfterBuild["Brick"])
// 		}
// 		if player1ResourcesAfterBuild["Sheep"] != 2 {
// 			t.Errorf("expected to have 1 Sheep after build city, but found %d", player1ResourcesAfterBuild["Sheep"])
// 		}
// 		if player1ResourcesAfterBuild["Grain"] != 1 {
// 			t.Errorf("expected to have 1 Grain after build city, but found %d", player1ResourcesAfterBuild["Grain"])
// 		}
// 		if player1ResourcesAfterBuild["Ore"] != 1 {
// 			t.Errorf("expected to have 1 Ore after build city, but found %d", player1ResourcesAfterBuild["Ore"])
// 		}
// 	})
// }

// func TestBuildCityErrorAlreadyExistsByPlayer(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithCitiesByPlayer(map[string][]int{
// 			"1": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  3,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build error - player has city in vertex", func(t *testing.T) {
// 		err := game.BuildCity("1", 42)
// 		if err == nil {
// 			t.Errorf("expected to not be able to build city in vertex#32, but it built just fine")
// 		}
// 	})
// }
//
// func TestBuildCityErrorSettlementAlreadyExistsOtherPlayer(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithSettlementsByPlayer(map[string][]int{
// 			"2": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  3,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build error - another player has settlement in vertex", func(t *testing.T) {
// 		err := game.BuildCity("1", 42)
// 		if err == nil {
// 			t.Errorf("expected to not be able to build city in vertex#32, but it built just fine")
// 		}
// 	})
// }
//
// func TestBuildCityErrorCityAlreadyExistsOtherPlayer(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithCitiesByPlayer(map[string][]int{
// 			"2": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  3,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build error - another player has city in vertex", func(t *testing.T) {
// 		err := game.BuildCity("1", 42)
// 		if err == nil {
// 			t.Errorf("expected to not be able to build city in vertex#32, but it built just fine")
// 		}
// 	})
// }
//
// func TestBuildCityErrorNotPlayerRound(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithCurrentRoundPlayer("2"),
// 		MockWithSettlementsByPlayer(map[string][]int{
// 			"1": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  3,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build error - it's not the player's round", func(t *testing.T) {
// 		err := game.BuildCity("1", 42)
// 		if err == nil {
// 			t.Errorf("expected to not be able to build city off round, but it built just fine")
// 		}
// 	})
// }
//
// func TestBuildCityErrorNotEnoughResources(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithSettlementsByPlayer(map[string][]int{
// 			"1": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  1,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build error - player doesn't have enough resources", func(t *testing.T) {
// 		err := game.BuildCity("1", 42)
// 		if err == nil {
// 			t.Errorf("expected to not be able to build city without enough resources, but it built just fine")
// 		}
// 	})
// }
//
// func TestBuildCityErrorNotAppropriateRound(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(MoveRobberDue7),
// 		MockWithSettlementsByPlayer(map[string][]int{
// 			"1": {42},
// 		}),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 0,
// 				"Brick":  1,
// 				"Sheep":  2,
// 				"Grain":  3,
// 				"Ore":    4,
// 			},
// 		}),
// 	)
//
// 	t.Run("city build error - player tries to build in an appropriate phase", func(t *testing.T) {
// 		err := game.BuildCity("1", 42)
// 		if err == nil {
// 			t.Errorf("expected to not be able to build city without being in regular phase, but it built just fine")
// 		}
// 	})
// }

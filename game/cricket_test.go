package game

import (
	"fmt"
	"testing"

	"github.com/gocaine/go-dart/common"
)

func TestRegular2Players(t *testing.T) {
	fmt.Println()
	fmt.Println("TestRegular2Players")

	ctx := createContext("eng")

	game, err := NewGameCricket(ctx, map[string]interface{}{"CutThroat": true, "NoScore": true})

	expected := "CutThroat and NoScore options are not compatible"
	AssertError(t, err, expected)

	game, err = NewGameCricket(ctx, map[string]interface{}{})

	err = game.Start(ctx)
	AssertError(t, err, "Game cannot start")

	game.AddPlayer(ctx, "test_board", "Alice")
	game.AddPlayer(ctx, "test_board", "Bob")

	state := game.State()

	_, err = game.HandleDart(ctx, common.Sector{Val: 21, Pos: 2})
	AssertError(t, err, "Sector is not a valid one")

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 8, Pos: 2})
	AssertCurrents(t, state, 0, 1)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 2})
	AssertCurrents(t, state, 0, 2)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 1})
	AssertCurrents(t, state, 0, 2)

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})
	assertHistoOrMemory(t, game.state.Players[1].Histo, "16", 3)
	assertHistoOrMemory(t, game.memory, "16", 1)
	AssertCurrents(t, state, 1, 1)
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 1})
	AssertScore(t, state.Players[1], 16)
	AssertCurrents(t, state, 1, 2)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 3})
	assertHistoOrMemory(t, game.memory, "15", 1)

	game.HoldOrNextPlayer(ctx)

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 3})
	AssertScore(t, state.Players[0], 0)
	AssertEquals(t, game.memory["15"], 0)
	assertHistoOrMemory(t, game.memory, "15", 0)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertScore(t, state.Players[0], 20)
	game.HandleDart(ctx, common.Sector{Val: 18, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 1})

	game.HoldOrNextPlayer(ctx)

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 18, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 2})
	assertHistoOrMemory(t, game.state.Players[1].Histo, "15", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "16", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "17", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "18", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "19", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "20", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "25", 3)

	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.PLAYING)

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})
	AssertScore(t, state.Players[0], 20)

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	AssertScore(t, state.Players[1], 16)
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 1})

	AssertGameState(t, state, common.OVER)

	player := state.Players[0]

	AssertScore(t, player, 41)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 20)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")

}

func TestNoScore2Players(t *testing.T) {
	fmt.Println()
	fmt.Println("TestNoScore2Players")

	ctx := createContext("eng")

	game, _ := NewGameCricket(ctx, map[string]interface{}{"NoScore": true})

	game.AddPlayer(ctx, "test_board", "Alice")
	game.AddPlayer(ctx, "test_board", "Bob")

	state := game.State()

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 8, Pos: 2})
	AssertCurrents(t, state, 0, 1)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 2})
	AssertCurrents(t, state, 0, 2)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 1})
	AssertCurrents(t, state, 0, 2)

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})
	assertHistoOrMemory(t, game.state.Players[1].Histo, "16", 3)
	assertHistoOrMemory(t, game.memory, "16", 1)
	AssertCurrents(t, state, 1, 1)
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 1})
	AssertScore(t, state.Players[1], 0)
	AssertCurrents(t, state, 1, 2)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 3})
	assertHistoOrMemory(t, game.memory, "15", 1)

	game.HoldOrNextPlayer(ctx)

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 3})
	AssertScore(t, state.Players[0], 0)
	AssertEquals(t, game.memory["15"], 0)
	assertHistoOrMemory(t, game.memory, "15", 0)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 18, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 1})

	game.HoldOrNextPlayer(ctx)

	// Alice
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Bob
	game.HandleDart(ctx, common.Sector{Val: 18, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 2})

	AssertGameState(t, state, common.OVER)

	player := state.Players[0]

	AssertScore(t, player, 0)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	assertHistoOrMemory(t, player.Histo, "15", 3)
	assertHistoOrMemory(t, player.Histo, "16", 3)
	assertHistoOrMemory(t, player.Histo, "17", 3)
	assertHistoOrMemory(t, player.Histo, "18", 3)
	assertHistoOrMemory(t, player.Histo, "19", 3)
	assertHistoOrMemory(t, player.Histo, "20", 3)
	assertHistoOrMemory(t, player.Histo, "25", 3)

	player = state.Players[1]

	AssertScore(t, player, 0)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")
}

func TestCutThroat3Players(t *testing.T) {
	fmt.Println()
	fmt.Println("TestCutThroat3Players")

	ctx := createContext("eng")

	game, _ := NewGameCricket(ctx, map[string]interface{}{"CutThroat": true})

	game.AddPlayer(ctx, "test_board", "Alice")
	game.AddPlayer(ctx, "test_board", "Bob")
	game.AddPlayer(ctx, "test_board", "Charly")

	state := game.State()

	// Alice 15:2 20:1
	AssertCurrents(t, state, 0, 0)
	game.HandleDart(ctx, common.Sector{Val: 8, Pos: 2})
	AssertCurrents(t, state, 0, 1)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 2})
	AssertCurrents(t, state, 0, 2)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 1})
	AssertCurrents(t, state, 0, 2)

	game.HoldOrNextPlayer(ctx)

	// Bob 16:3 15:3
	AssertCurrents(t, state, 1, 0)
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})
	assertHistoOrMemory(t, game.state.Players[1].Histo, "16", 3)
	assertHistoOrMemory(t, game.memory, "16", 2)
	AssertCurrents(t, state, 1, 1)
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 1})
	AssertScore(t, state.Players[1], 0)
	AssertScore(t, state.Players[0], 16)
	AssertScore(t, state.Players[2], 16)
	AssertCurrents(t, state, 1, 2)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 3})
	assertHistoOrMemory(t, game.memory, "15", 2)

	game.HoldOrNextPlayer(ctx)

	// Charly
	AssertCurrents(t, state, 2, 0)
	game.HandleDart(ctx, common.Sector{Val: 7, Pos: 2})
	game.HandleDart(ctx, common.Sector{Val: 11, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 3, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Alice 15:3 20:3 18:3
	AssertCurrents(t, state, 0, 0)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 3})
	AssertScore(t, state.Players[0], 16)
	AssertScore(t, state.Players[1], 0)
	AssertScore(t, state.Players[2], 46)
	assertHistoOrMemory(t, game.memory, "15", 1)
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	AssertScore(t, state.Players[0], 16)
	AssertScore(t, state.Players[1], 20)
	AssertScore(t, state.Players[2], 66)
	game.HandleDart(ctx, common.Sector{Val: 18, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Bob 15:3 16:3 17:3 19:3 25:1
	AssertCurrents(t, state, 1, 0)
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 1})

	game.HoldOrNextPlayer(ctx)

	// Charly 15:3
	AssertCurrents(t, state, 2, 0)
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 15, Pos: 2})
	assertHistoOrMemory(t, game.memory, "15", 0)

	game.HoldOrNextPlayer(ctx)

	// Alice 15:3 16:3 17:3 18:3 19:3 20:3
	AssertCurrents(t, state, 0, 0)
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Bob 15:3 16:3 17:3 18:3 19:3 20:3 25:3
	AssertCurrents(t, state, 1, 0)
	game.HandleDart(ctx, common.Sector{Val: 18, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 20, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 2})
	assertHistoOrMemory(t, game.state.Players[1].Histo, "15", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "16", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "17", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "18", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "19", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "20", 3)
	assertHistoOrMemory(t, game.state.Players[1].Histo, "25", 3)

	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.PLAYING)

	// Charly 15:3
	AssertCurrents(t, state, 2, 0)
	game.HandleDart(ctx, common.Sector{Val: 7, Pos: 2})
	game.HandleDart(ctx, common.Sector{Val: 11, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 3, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Alice 15:3 16:3 17:3 18:3 19:3 20:3
	AssertCurrents(t, state, 0, 0)
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 3})
	game.HandleDart(ctx, common.Sector{Val: 16, Pos: 3})
	AssertScore(t, state.Players[0], 16)
	AssertScore(t, state.Players[1], 20)
	AssertScore(t, state.Players[2], 222)

	game.HoldOrNextPlayer(ctx)

	// Bob 15:3 16:3 17:3 18:3 19:3 20:3 25:3
	AssertCurrents(t, state, 1, 0)
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 3})
	AssertScore(t, state.Players[0], 16)
	AssertScore(t, state.Players[1], 20)
	AssertScore(t, state.Players[2], 279)
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 1})
	AssertRank(t, state.Players[1], 1)
	AssertScore(t, state.Players[0], 41)
	AssertScore(t, state.Players[1], 20)
	AssertScore(t, state.Players[2], 304)

	game.HoldOrNextPlayer(ctx)

	// Charly 15:3
	AssertCurrents(t, state, 2, 0)
	game.HandleDart(ctx, common.Sector{Val: 7, Pos: 2})
	game.HandleDart(ctx, common.Sector{Val: 11, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 3, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Alice 15:3 16:3 17:3 18:3 19:3 20:3 25:1
	AssertCurrents(t, state, 0, 0)
	game.HandleDart(ctx, common.Sector{Val: 19, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 17, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 1})
	AssertScore(t, state.Players[0], 41)
	AssertScore(t, state.Players[1], 20)
	AssertScore(t, state.Players[2], 340)

	game.HoldOrNextPlayer(ctx)

	// Charly 15:3
	AssertCurrents(t, state, 2, 0)
	game.HandleDart(ctx, common.Sector{Val: 7, Pos: 2})
	game.HandleDart(ctx, common.Sector{Val: 11, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 3, Pos: 3})

	game.HoldOrNextPlayer(ctx)

	// Alice 15:3 16:3 17:3 18:3 19:3 20:3 25:3
	AssertCurrents(t, state, 0, 0)
	game.HandleDart(ctx, common.Sector{Val: 25, Pos: 2})

	AssertGameState(t, state, common.OVER)

	player := state.Players[0]

	AssertScore(t, player, 20)
	AssertRank(t, player, 1)
	AssertName(t, player, "Bob")

	player = state.Players[1]

	AssertScore(t, player, 41)
	AssertRank(t, player, 2)
	AssertName(t, player, "Alice")

	player = state.Players[2]

	AssertScore(t, player, 340)
	AssertRank(t, player, 3)
	AssertName(t, player, "Charly")

}

func TestGameCricketOnHold(t *testing.T) {
	fmt.Println()
	fmt.Println("TestGameCricketOnHold")

	ctx := createContext("eng")

	game, _ := NewGameCricket(ctx, map[string]interface{}{"NoScore": true})
	game.AddPlayer(ctx, "test_board", "Alice")
	game.AddPlayer(ctx, "test_board", "Bob")

	state, _ := game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 0, 1)

	_, err := game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	AssertError(t, err, "Game is on hold and not ready to handle darts")

	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.PLAYING)
	AssertCurrents(t, state, 1, 0)

	game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})
	game.HandleDart(ctx, common.Sector{Val: 5, Pos: 1})

	AssertGameState(t, state, common.ONHOLD)
	AssertCurrents(t, state, 1, 2)

	game.HoldOrNextPlayer(ctx)

	AssertGameState(t, state, common.PLAYING)
	AssertCurrents(t, state, 0, 0)
}

func assertHistoOrMemory(t *testing.T, histo map[string]int, key string, val int) {
	if histo[key] != val {
		fatalStack(t, "HistoOrMemory : Map[%s] should be %d but was %d -- %+v", key, val, histo[key], histo)
	}
}

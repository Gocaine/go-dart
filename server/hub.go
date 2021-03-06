package server

import (
	"encoding/json"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocaine/go-dart/common"
	"github.com/gocaine/go-dart/game"
	"github.com/gocaine/go-dart/i18n"
	"golang.org/x/net/websocket"
)

// GameHub handle websocket connections for a Game
type GameHub struct {
	clients []hubClient
	output  chan bool
	game    game.Game
}

type hubClient struct {
	locale string
	ws     *websocket.Conn
}

type wsMessage struct {
	Kind string
	Data interface{}
}

// NewGameHub is GameHub constructor
func NewGameHub(game game.Game) *GameHub {
	hub := GameHub{game: game, clients: make([]hubClient, 0)}
	return &hub
}

func (gh *GameHub) handle(locale string) func(*websocket.Conn) {
	return func(connection *websocket.Conn) {
		log.Infof("new ws connection for this user")
		gh.clients = append(gh.clients, hubClient{ws: connection, locale: locale})
		statusAsBytes := dataAsBytes("status", gh.game.State())
		connection.Write([]byte(statusAsBytes))
		// lock until the end of the world
		connection.Read(make([]byte, 0))
	}
}

func (gh *GameHub) refresh() {
	statusAsBytes := dataAsBytes("status", gh.game.State())
	for _, client := range gh.clients {
		log.Info("sending status")
		_, err := client.ws.Write(statusAsBytes)
		if err != nil {
			log.Infof("error writing %v", err)
		}
	}
}

func (gh *GameHub) close() {
	log.Infof("close all websocket connections")
	for _, client := range gh.clients {
		client.ws.Close()
	}
}

func (gh *GameHub) sendMessage(key string, args ...interface{}) {
	for _, client := range gh.clients {
		log.Info("sending message")
		msgAsBytes := dataAsBytes("message", fmt.Sprintf(i18n.Translation(key, client.locale), args...))
		_, err := client.ws.Write(msgAsBytes)
		if err != nil {
			log.Infof("error writing %v", err)
		}
	}
}

func dataAsBytes(kind string, data interface{}) []byte {
	dataAsBytes, err := json.Marshal(wsMessage{Kind: kind, Data: data})
	if err != nil {
		log.Infof("cannot serialize data : %v", data)
	}
	return dataAsBytes
}

// Start start the game, Darts will be handled
func (gh *GameHub) Start(ctx common.GameContext) error {
	ctx.MessageHandler = gh.sendMessage
	return gh.game.Start(ctx)
}

// AddPlayer add a new player to the game
func (gh *GameHub) AddPlayer(ctx common.GameContext, board string, name string) error {
	ctx.MessageHandler = gh.sendMessage
	return gh.game.AddPlayer(ctx, board, name)
}

// HandleDart the implementation has to handle the Dart regarding the current player, the rules, and the context. Return a GameState
func (gh *GameHub) HandleDart(ctx common.GameContext, sector common.Sector) (*common.GameState, error) {
	ctx.MessageHandler = gh.sendMessage
	return gh.game.HandleDart(ctx, sector)
}

// State returns the current GameState
func (gh *GameHub) State() *common.GameState {
	return gh.game.State()
}

// BoardHasLeft is call to notify the game a board has been disconnected. Returns true if the game continues despite this event
func (gh *GameHub) BoardHasLeft(ctx common.GameContext, board string) bool {
	ctx.MessageHandler = gh.sendMessage
	return gh.game.BoardHasLeft(ctx, board)
}

// HoldOrNextPlayer switch game state between ONHOLD and PLAYING with side effects according to game implementation
func (gh *GameHub) HoldOrNextPlayer(ctx common.GameContext) {
	ctx.MessageHandler = gh.sendMessage
	gh.game.HoldOrNextPlayer(ctx)
}

// NextPlayer is called when the current player end his visit
func (gh *GameHub) NextPlayer(ctx common.GameContext) {
	ctx.MessageHandler = gh.sendMessage
	gh.game.NextPlayer(ctx)
}

// NextDart is called after each dart when the same palyer play again
func (gh *GameHub) NextDart(ctx common.GameContext) {
	ctx.MessageHandler = gh.sendMessage
	gh.game.NextDart(ctx)
}

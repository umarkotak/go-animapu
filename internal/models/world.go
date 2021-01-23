package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type WorldMap struct {
}

type WorldInfo struct {
	WorldMap         WorldMap
	HorizontalLength int64
	VerticalLength   int64
}

type World struct {
	Players       map[*Player]bool
	LoggedPlayers map[*Player]bool
	PlayerDB      map[string]*Player
	Broadcast     chan []byte
	Register      chan *Player
	Unregister    chan *Player
	WorldInfo     WorldInfo
}

func NewWorldInfo() WorldInfo {
	return WorldInfo{
		WorldMap:         WorldMap{},
		HorizontalLength: 10,
		VerticalLength:   10,
	}
}

func NewWorld() *World {
	return &World{
		Players:       make(map[*Player]bool),
		LoggedPlayers: make(map[*Player]bool),
		PlayerDB:      make(map[string]*Player),
		Broadcast:     make(chan []byte),
		Register:      make(chan *Player),
		Unregister:    make(chan *Player),
		WorldInfo:     NewWorldInfo(),
	}
}

func (w *World) Run() {
	for {
		select {
		case player := <-w.Register:
			w.Players[player] = true
			w.PlayerDB[player.Name] = player

		case player := <-w.Unregister:
			if _, ok := w.Players[player]; ok {
				delete(w.Players, player)
				close(player.Send)
			}

		case message := <-w.Broadcast:
			var gameMessage GameMessage
			json.Unmarshal(message, &gameMessage)
			w.handleGameMessage(gameMessage)
		}
	}
}

func (w *World) handleGameMessage(gameMessage GameMessage) {
	switch gameMessage.MessageType {
	case GmGlobalMessage:
		w.serviceGlobalMessage(gameMessage)
	case GmWorldMapInfo:
		w.serviceWorldMapInfo(gameMessage)
	default:
	}
}

func (w *World) serviceGlobalMessage(gameMessage GameMessage) {
	message := fmt.Sprintf("%v", gameMessage.Data["message"])

	responseBroadcast := GameMessage{
		MessageType: GmGlobalMessage,
		Data: map[string]interface{}{
			"from":    gameMessage.Player.Name,
			"message": message,
			"ts":      time.Now().UTC(),
		},
		Direction: "response",
	}
	response, _ := json.Marshal(responseBroadcast)

	for player := range w.Players {
		select {
		case player.Send <- []byte(response):
		default:
			close(player.Send)
			delete(w.Players, player)
		}
	}
}

func (w *World) serviceWorldMapInfo(gameMessage GameMessage) {
	selectedPlayer := w.PlayerDB[gameMessage.Player.Name]
	tempPlayers := map[string][]string{}

	for tempPlayer := range w.Players {
		key := fmt.Sprintf("%v-%v", tempPlayer.PlayerInfo.PosY, tempPlayer.PlayerInfo.PosX)
		tempPlayers[key] = append(tempPlayers[key], tempPlayer.Name)
	}

	maps := []interface{}{}

	for y := int64(1); y <= w.WorldInfo.VerticalLength; y++ {
		row := []interface{}{}
		for x := int64(1); x <= w.WorldInfo.HorizontalLength; x++ {
			rowKey := fmt.Sprintf("%v-%v", y, x)
			rowData := map[string]interface{}{
				"pos_x": x,
				"pos_y": y,
				"info": map[string]interface{}{
					"terrain": "grassland",
				},
				"players": tempPlayers[rowKey],
			}

			row = append(row, rowData)
		}

		maps = append(maps, row)
	}

	responseBroadcast := GameMessage{
		MessageType: GmWorldMapInfo,
		Data: map[string]interface{}{
			"horizontal_length": w.WorldInfo.HorizontalLength,
			"vertical_length":   w.WorldInfo.VerticalLength,
			"maps":              maps,
		},
		Direction: "response",
	}
	response, _ := json.Marshal(responseBroadcast)

	selectedPlayer.Send <- []byte(response)
}

func (w *World) serviceError(gameMessage GameMessage) {

}

package models

const (
	GmInitStream    = "init_stream"
	GmPlayerInfo    = "player_info"
	GmWorldMapInfo  = "world_map_info"
	GmGlobalMessage = "global_message"
)

type GameMessage struct {
	Player      Player                 `json:"player"`
	MessageType string                 `json:"message_type"`
	Meta        map[string]interface{} `json:"meta"`
	Data        map[string]interface{} `json:"data"`
	Headers     GameHeader             `json:"headers"`
	Direction   string                 `json:"direction"`
}

type GameHeader struct {
	Authorization string `json:"authorization"`
}

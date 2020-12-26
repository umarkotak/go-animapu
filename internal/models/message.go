package models

type Message struct {
	MessageType   string `json:"message_type"`
	Name          string `json:"name"`
	Message       string `json:"message"`
	TimestampUnix int32  `json:"timestamp_unix"`
	TimestampUTC  string `json:"timestamp_utc"`
}

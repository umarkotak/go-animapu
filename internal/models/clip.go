package models

type Clip struct {
	ID            string `json:"id"`
	Content       string `json:"content"`
	TimestampUnix int32  `json:"timestamp_unix"`
	TimestampUtc  string `json:"timestamp_utc"`
}

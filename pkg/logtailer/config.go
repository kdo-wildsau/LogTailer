package logtailer

// Config logTailer Config file definition
type Config struct {
	LogPath           string `json:"logPath"`
	LogPathUpdateTime int    `json:"logPathUpdateTime"`
	MemzoInstanceName string `json:"memzoInstanceName"`
	MemzoIngestionKey string `json:"memzoIngestionKey"`
}

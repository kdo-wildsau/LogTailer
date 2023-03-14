package logtailer

type Config struct {
	LogPath           string `json:"logPath"`
	LogPathUpdateTime int    `json:"LogPathUpdateTime"`
	DiscordWebhookUrl string `json:"discordWebhookUrl"`
}

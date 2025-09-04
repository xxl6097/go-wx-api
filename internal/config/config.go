package config

type Config struct {
	ServerPort int    `json:"serverPort"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AppID      string `json:"appid"`
	AppSecret  string `json:"appsecret"`
}

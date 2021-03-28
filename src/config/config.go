package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/stevenhay/qiqi_bot/src/io"
)

type Config struct {
	GuildID        string      `json:"guild_id"`
	GuildOwnerID   string      `json:"guild_owner_id"`
	BotOwnerID     string      `json:"bot_owner_id"`
	AuditChannelID string      `json:"audit_channel_id"`
	ErrorChannelID string      `json:"error_channel_id"`
	DefaultRole    DefaultRole `json:"default_role"`
}

type DefaultRole struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

var _config *Config

func Get() *Config {
	return _config
}

func (c Config) Save() error {
	return io.SaveConfig("guild", c)
}

func Load() {
	mappingFile, err := os.Open("../config/guild.json")
	if err != nil {
		logrus.Fatal("Unable to open guild config file:", err)
	}
	defer mappingFile.Close()

	byteVal, err := ioutil.ReadAll(mappingFile)
	if err != nil {
		logrus.Fatal("Unable to read guild config file:", err)
	}

	var res Config
	json.Unmarshal([]byte(byteVal), &res)

	fmt.Println(res)
	_config = &res
}

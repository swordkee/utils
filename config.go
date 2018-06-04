package utils

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Servers    map[string]servers
	DateBase   map[string]dateBase
	Futures     map[string]futures
	Adapter    map[string]adapter
	Indicators map[string]indicators
	Alert      map[string]alert
	Msg        map[string]msgCode
	FBDebug    bool
}

type dateBase struct {
	Host    string
	Port    int
	User    string
	Pass    string
	DbName  string
	Charset string
	Pool    int
}

type servers struct {
	Port int
	Host string
	Path string
}
type futures struct {
	BrokerID    string
	MarketFront string
	TradeFront  string
}
type adapter struct {
	Host     string
	Port     int
	Protocol string
	Path     string
	Symbol   []string
}
type alert struct {
	KTime      int
	Symbols    []string
	Indicators []string
}
type indicators struct {
	Spacing      int
	CloseSpacing int
	AtrParam     int
	AtrMultiple  int
}
type msgCode struct {
	SignName        string
	GatewayUrl      string
	AccessKeyId     string
	AccessKeySecret string
	TmplCode        string
}

func GetConfig() TomlConfig {
	debug := BoolMust(os.Getenv("FB_DEBUG"))
	var config TomlConfig
	config.FBDebug = debug
	if _, err := toml.DecodeFile(GetConfigPath(), &config); err != nil {
		log.Println(err)
		return config
	}
	return config
}

func GetConfigPath() string {
	debug := BoolMust(os.Getenv("FB_DEBUG"))
	tomlFile := On(!debug, "config", "config_debug").(string)
	return SelfDir() + "/config/" + tomlFile + ".toml"
}

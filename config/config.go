package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type GlobalConfig struct {
	ListenServer  string
	ListenPort    int
	LogPath       string
	LogReserveDay int
	StaticDir     string
	Domain        string
}

type DBConfig struct {
	User     string
	Pass     string
	Host     string
	Port     int
	Database string
}

type RedisConfig struct {
	Host     string
	Port     int
	Pass     string
	Database int
	PoolSize int
}

type LaunchModule struct {
	// oss 有很多功能模块，这里配置这模块的开关
	LaunchApi bool
}

type Config struct {
	GlobalConfig
	DB           DBConfig     `toml:"db"`
	RedisConfig  RedisConfig  `toml:"redis"`
	LaunchModule LaunchModule `toml:"launch"`
}

func DefaultConfig() Config {
	baseDir := filepath.Dir(os.Args[0])
	return Config{
		GlobalConfig: GlobalConfig{
			ListenServer:  "",
			ListenPort:    80,
			LogPath:       filepath.Join(baseDir, "../log/sky.log"),
			LogReserveDay: 90,
		},
	}
}

func NewConfig(configFile string) (*Config, error) {
	config := DefaultConfig()
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	dir, _ := os.Getwd()
	configDir := fmt.Sprintf("%s/config", dir)
	log.Printf("config dir: %s", configDir)
	config.AddConfigPath(configDir)
	// config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatalf("error on parsing configuration file, %v", err)
	}
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

// GetConfig ..
func GetConfig() *viper.Viper {
	return config
}

package logger

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

var config *viper.Viper

func init() {
	godotenv.Load()
	cmd, _ := os.Getwd()
	config = viper.New()
	config.SetConfigName(os.Getenv("ENV"))
	config.SetConfigType("yaml")
	config.AddConfigPath(path.Join(cmd, "config"))
	err := config.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("fatal error config file: %w", err))
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
}

func Get(key string) any {
	return config.Get(key)
}

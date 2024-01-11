package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

var viperConfig *viper.Viper

func init() {
	godotenv.Load()
	cmd, _ := os.Getwd()
	viperConfig = viper.New()
	viperConfig.SetConfigName(fmt.Sprintf("config_%v", os.Getenv("ENV")))
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath(path.Join(cmd, "config"))
	err := viperConfig.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("fatal error config file: %v", err))
	}
}

// 获取字符串类型的配置
func GetString(key string) string {
	return viperConfig.GetString(key)
}

// 获取布尔类型的配置
func GetBool(key string) bool {
	return viperConfig.GetBool(key)
}

// 获取整型的配置
func GetInt(key string) int {
	return viperConfig.GetInt(key)
}

package bootstrap

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var viperConfig *viper.Viper

type config struct{}

var Config = newConfig()
var logger = Logger.GetLogger()

func newConfig() *config {
	godotenv.Load()
	cmd, _ := os.Getwd()
	viperConfig = viper.New()
	viperConfig.SetConfigName(os.Getenv("ENV"))
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath(path.Join(cmd, "config"))
	err := viperConfig.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("fatal error config file: %w", err))
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return &config{}
}

// 获取布尔类型的配置
func (c *config) GetBoolValue(key string) bool {
	value, err := strconv.ParseBool(fmt.Sprintf("%s", viperConfig.Get(key)))
	if err != nil {
		logger.DPanic(err.Error())
	}
	return value
}

// 获取字符串类型的配置
func (c *config) GetStringValue(key string) string {
	value := viperConfig.Get(key)
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", viperConfig.Get(key)))
}

// 获取整型的配置
func (c *config) GetIntegerValue(key string) int {
	value, err := strconv.Atoi(c.GetStringValue(key))
	if err != nil {
		logger.DPanic(err.Error())
	}
	return value
}

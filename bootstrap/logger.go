package bootstrap

import (
	"fmt"
	"os"
)

func init() {
	cmd, _ := os.Getwd()
	fmt.Println("bootstrap logger", cmd)
}

func NewLogger() {
	fmt.Println("New xxxxxxxxx Logger")
}

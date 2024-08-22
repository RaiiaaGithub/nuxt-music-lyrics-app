package utils

import (
	"fmt"
	"time"
)

func getFullDate() string {
	now := time.Now()
	return fmt.Sprintf("%d/%d/%d %02d:%02d:%02d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
}

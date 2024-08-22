package utils

import "fmt"

var BASE_MESSAGE = getFullDate() + " [RAIA] - "

var dev bool = true

func LogDebug(title string, value ...any) {
	if dev {
		fmt.Printf(BASE_MESSAGE+title+"\n", value...)
	}
}

func LogError(title string, value ...any) {
	if dev {
		fmt.Printf(BASE_MESSAGE+title+"\n", value...)
	}
}

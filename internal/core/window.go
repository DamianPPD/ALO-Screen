package core

import "strings"

func IsChatGPTWindow(title string) bool {
	return strings.Contains(strings.ToLower(title), "chatgpt")
}

package main

import (
	"chargpt_ui/chatgpt_ui"
	"fmt"
)

func main() {
	err := chatgpt_ui.ChatGPTUi()
	fmt.Println("err: ", err)
}

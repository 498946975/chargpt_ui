package main

import (
	"chargpt_ui/chartgpt_ui"
	"fmt"
)

func main() {
	err := chartgpt_ui.ChartGPTUi()
	fmt.Println("err: ", err)
}

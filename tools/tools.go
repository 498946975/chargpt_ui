package tools

import (
	"encoding/json"
	"github.com/andlabs/ui"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ApiKey      string `json:"ApiKey"`
	HttpTimeOut int    `json:"HttpTimeOut"`
}

func MyEntry(init string) *ui.Entry {
	entry := ui.NewEntry()
	entry.SetText(init)
	entry.LibuiControl()
	return entry
}

func MyGroup(group_name string, entry *ui.Entry) *ui.Group {
	group := ui.NewGroup(group_name)
	group.SetChild(entry)
	return group
}

func MyBox(group *ui.Group) *ui.Box {
	boxs := ui.NewHorizontalBox()
	boxs.Append(group, true)
	boxs.SetPadded(false)
	return boxs
}

func StringToFloat32(str string) float32 {
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		log.Println("Error:", err)
	}
	f32 := float32(f)
	return f32
}

func StringToInt(str string) int {
	n, err := strconv.Atoi(str)
	if err != nil {
		log.Println("Error:", err)
	}
	return n
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

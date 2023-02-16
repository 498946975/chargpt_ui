package chatgpt_ui

import (
	"bytes"
	"chargpt_ui/tools"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"
)

type CompletionResponse struct {
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
}

type CompletionRequest struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"top_p"`
	PresencePenalty  float32 `json:"presence_penalty"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
}

func ChatGPTHttp(model string, completionRequest CompletionRequest, cfg *tools.Config) string {
	var completionResponse CompletionResponse
	requestBody, _ := json.Marshal(completionRequest)
	url := fmt.Sprintf("https://api.openai.com/v1/engines/%s/completions", model)
	apiKey := fmt.Sprintf("Bearer %s", cfg.ApiKey)
	client := &http.Client{
		Timeout: time.Duration(cfg.HttpTimeOut) * time.Second,
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)

	res, http_err := client.Do(req)
	if http_err != nil {
		return "TimeOut"
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "json decode failed"
	}
	return completionResponse.Choices[0].Text

}

func ChatGPTUi() error {
	var completion_request CompletionRequest
	relativePath, _ := filepath.Abs("./cfg.json")
	cfg, err := tools.LoadConfig(relativePath)
	if err != nil {
		panic(err)
	}
	logger, _ := tools.NewLogger(cfg.LogDir)
	models := []string{"text-davinci-003", "code-davinci-002"}
	err = ui.Main(func() {
		prompt := tools.MyEntry("")
		max_tokens := tools.MyEntry("4000")
		temperature := tools.MyEntry("0.9")
		top_p := tools.MyEntry("1.0")
		presence_penalty := tools.MyEntry("0.0")
		frequency_penalty := tools.MyEntry("0.0")
		http_result_txt := ui.NewLabel("")

		model := ui.NewRadioButtons()
		model.Append(models[0])
		model.Append(models[1])
		model.SetSelected(0)

		parseButton := ui.NewButton("Parse")
		parseButton.OnClicked(func(*ui.Button) {
			out, err := exec.Command("pbpaste").Output()
			if err != nil {
				fmt.Println(err)
			}
			prompt.SetText(string(out))
		})
		copyButton := ui.NewButton("Copy")
		copyButton.OnClicked(func(*ui.Button) {
			text := http_result_txt.Text()
			cmd := exec.Command("pbcopy")
			in, _ := cmd.StdinPipe()
			_, err = in.Write([]byte(text))
			err = in.Close()
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		})

		commitButton := ui.NewButton("Submit")
		commitButton.OnClicked(func(*ui.Button) {
			completion_request.Prompt = prompt.Text()
			logger.Log("INFO", prompt.Text())
			completion_request.MaxTokens = tools.StringToInt(max_tokens.Text())
			completion_request.Temperature = tools.StringToFloat32(temperature.Text())
			completion_request.TopP = tools.StringToFloat32(top_p.Text())
			completion_request.PresencePenalty = tools.StringToFloat32(presence_penalty.Text())
			completion_request.FrequencyPenalty = tools.StringToFloat32(frequency_penalty.Text())
			resultCh := make(chan string)
			go func() {
				resultCh <- ChatGPTHttp(models[model.Selected()], completion_request, cfg)
			}()
			select {
			case result := <-resultCh:
				http_result_txt.SetText(result)
				logger.Log("WARNING", result)
			}

		})

		promptGroup := tools.MyGroup("Prompt", prompt)
		maxTokensGroup := tools.MyGroup("MaxTokens：0~4000", max_tokens)
		temperatureGroup := tools.MyGroup("Temperature：0.0~0.9", temperature)
		topPGroup := tools.MyGroup("TopP：0.1~1.0", top_p)
		presencePenaltyGroup := tools.MyGroup("PresencePenalty：-2.0~2.0", presence_penalty)
		frequencyPenaltyGroup := tools.MyGroup("FrequencyPenalty：-2.0~2.0", frequency_penalty)

		buttonGroup := ui.NewGroup("Submit request parameters")
		buttonGroup.SetChild(commitButton)

		copyGroup := ui.NewGroup("Copy the result of ChatGPT")
		copyGroup.SetChild(copyButton)

		parseGroup := ui.NewGroup("Paste the query")
		parseGroup.SetChild(parseButton)

		modelGroup := ui.NewGroup("Models")
		modelGroup.SetChild(model)

		resultGroup := ui.NewGroup("ChatGPT response")
		resultGroup.SetChild(http_result_txt)

		verticalBox := ui.NewVerticalBox()
		promptBox := tools.MyBox(promptGroup)
		parseBox := tools.MyBox(parseGroup)
		modelBox := tools.MyBox(modelGroup)
		maxTokensBox := tools.MyBox(maxTokensGroup)
		resultBox := tools.MyBox(resultGroup)

		boxsLine4 := ui.NewHorizontalBox()
		boxsLine4.Append(temperatureGroup, true)
		boxsLine4.Append(topPGroup, true)
		boxsLine4.SetPadded(false)

		boxsLine5 := ui.NewHorizontalBox()
		boxsLine5.Append(presencePenaltyGroup, true)
		boxsLine5.Append(frequencyPenaltyGroup, true)
		boxsLine5.SetPadded(false)

		boxsLine6 := ui.NewHorizontalBox()
		boxsLine6.Append(buttonGroup, true)
		boxsLine6.Append(copyGroup, true)
		boxsLine6.SetPadded(false)

		verticalBox.Append(promptBox, true)
		verticalBox.Append(parseBox, true)
		verticalBox.Append(modelBox, true)
		verticalBox.Append(maxTokensBox, true)
		verticalBox.Append(boxsLine4, true)
		verticalBox.Append(boxsLine5, true)
		verticalBox.Append(boxsLine6, true)
		verticalBox.Append(resultBox, false)
		verticalBox.SetPadded(false)

		window := ui.NewWindow("ChatGPT", 600, 600, false)
		window.SetChild(verticalBox)
		window.SetMargined(true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
	return err
}

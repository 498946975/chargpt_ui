package chartgpt_ui

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

func ChartGPTHttp(model string, completionRequest CompletionRequest) string {
	var completionResponse CompletionResponse
	relativePath, _ := filepath.Abs("./cfg.json")
	cfg, err := tools.LoadConfig(relativePath)
	if err != nil {
		panic(err)
	}
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

	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "json decode failed"
	}
	return completionResponse.Choices[0].Text

}

func ChartGPTUi() error {
	var completion_request CompletionRequest
	models := []string{"text-davinci-003", "code-davinci-002"}
	err := ui.Main(func() {
		prompt := tools.MyEntry("")
		max_tokens := tools.MyEntry("4000")
		temperature := tools.MyEntry("0.9")
		top_p := tools.MyEntry("1.0")
		presence_penalty := tools.MyEntry("0.0")
		frequency_penalty := tools.MyEntry("0.0")
		http_result_txt := tools.MyEntry("")

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
			_, err := in.Write([]byte(text))
			err = in.Close()
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		})

		commitButton := ui.NewButton("提交")
		commitButton.OnClicked(func(*ui.Button) {
			completion_request.Prompt = prompt.Text()
			completion_request.MaxTokens = tools.StringToInt(max_tokens.Text())
			completion_request.Temperature = tools.StringToFloat32(temperature.Text())
			completion_request.TopP = tools.StringToFloat32(top_p.Text())
			completion_request.PresencePenalty = tools.StringToFloat32(presence_penalty.Text())
			completion_request.FrequencyPenalty = tools.StringToFloat32(frequency_penalty.Text())
			resultCh := make(chan string)
			go func() {
				resultCh <- ChartGPTHttp(models[model.Selected()], completion_request)
			}()
			select {
			case result := <-resultCh:
				http_result_txt.SetText(result)
			}

		})

		promptGroup := tools.MyGroup("请输入需要获取的内容", prompt)
		maxTokensGroup := tools.MyGroup("设置返回信息的最大长度：0~4000", max_tokens)
		temperatureGroup := tools.MyGroup("随机性：0.0~0.9", temperature)
		topPGroup := tools.MyGroup("top_p：0.1~1.0", top_p)
		presencePenaltyGroup := tools.MyGroup("控制主题的重复度：-2.0~2.0", presence_penalty)
		frequencyPenaltyGroup := tools.MyGroup("控制字符的重复度：-2.0~2.0", frequency_penalty)

		buttonGroup := ui.NewGroup("提交请求参数")
		buttonGroup.SetChild(commitButton)

		copyGroup := ui.NewGroup("拷贝ChartGPT返回的结果")
		copyGroup.SetChild(copyButton)

		parseGroup := ui.NewGroup("复制需要查询的内容")
		parseGroup.SetChild(parseButton)

		modelGroup := ui.NewGroup("选择模型")
		modelGroup.SetChild(model)

		resultGroup := ui.NewGroup("ChartGPT返回信息")
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

		window := ui.NewWindow("ChartGPT", 600, 600, false)
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

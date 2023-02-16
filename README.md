## chartgpt使用
### 1、首先在cfg.json中修改apikey
```yaml
获取apikey: https://platform.openai.com/account/api-keys
```
### 2、MAC编译
```yaml
chargpt_ui目录下执行命令: go build -o ChatGPT main.go
```
### 3、Windows编译
```yaml
chargpt_ui目录下执行命令: CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o ChatGPT.exe main.go
```
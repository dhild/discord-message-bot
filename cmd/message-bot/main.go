package main

import (
	"flag"
	"os"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"encoding/json"
	"net/http"
	"bytes"
	"strings"
)

var (
	configFile string
	tts bool
)

func init() {
	envConfigFile := os.Getenv("MESSAGE_BOT_CONFIG_FILE")
	if envConfigFile == "" {
		envConfigFile = "message-bot.yaml"
	}
	flag.StringVar(&configFile, "config", envConfigFile, "override default config file (MESSAGE_BOT_CONFIG_FILE)")
	flag.BoolVar(&tts, "tts", false, "use Text-To-Speech")
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("no message given")
	}

	message := strings.Join(flag.Args(), " ")

	webhook, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}

	webhook.Execute(WebhookData{message, tts})
}

type Config struct {
	WebhookURL string `yaml:"webhookURL"`
}

func NewConfig(filename string) (*Config, error) {
	configData, err :=ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type WebhookData struct {
	Content string `json:"content"`
	TTS     bool `json:"tts"`
}

func (c *Config) Execute(webhookData WebhookData) (err error) {
	data, err := json.Marshal(webhookData)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(data)
	resp, err := http.Post(c.WebhookURL, "application/json", buf)
	defer resp.Body.Close()
	log.Printf("status: [%s], content length: [%v]", resp.Status, resp.ContentLength)
	return
}

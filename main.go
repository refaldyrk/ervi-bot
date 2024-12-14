package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/telebot.v3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result"`
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

var chatHistories = make(map[int64][]Message)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	botToken := viper.GetString("BOT_TOKEN")
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID
		userInput := c.Text()

		if c.Chat().Type == telebot.ChatGroup {
			isMentioned := false

			if strings.Contains(c.Text(), "@"+bot.Me.Username) {
				isMentioned = true
			}

			if !isMentioned && c.Message().ReplyTo == nil {

				return nil
			}

			fmt.Println("ON GROUP", isMentioned, c.Message().ReplyTo)
		}

		if _, exists := chatHistories[userID]; !exists {
			dataUserString, err := json.Marshal(c.Sender())
			if err != nil {
				return err
			}

			println("NEW", c.Sender().FirstName, c.Sender().LastName)
			chatHistories[userID] = []Message{
				{
					Role:    "system",
					Content: "Kamu adalah Ervi, seorang teman virtual. Tugasmu adalah merespons dengan gaya santai, asik, dan mengikuti gaya percakapan pengguna dan jangan baku santai aja. Kalau pengguna cuek, kamu bales cuek. Kalau marah, bales marah lagi. Kalau dia bercanda, ikut bercanda. Jadilah teman yang baik dan relatable. Berikut data pengguna untuk kamu kenali dalam bentuk json, anda harus kenal siapa yang mengobrol dengan anda: " + string(dataUserString),
				},
			}

		}

		chatHistories[userID] = append(chatHistories[userID], Message{
			Role:    "user",
			Content: userInput,
		})

		response, err := getAIResponse(chatHistories[userID])
		if err != nil {
			log.Printf("Error getting AI response: %v", err)
			return c.Reply("Waduh, aku lagi error nih. Coba lagi nanti ya~")
		}

		chatHistories[userID] = append(chatHistories[userID], Message{
			Role:    "assistant",
			Content: response,
		})

		return c.Reply(response)
	})

	fmt.Println("Bot Telegram berjalan...")
	bot.Start()
}

func getAIResponse(chatHistory []Message) (string, error) {
	apiURL := "https://api.cloudflare.com/client/v4/accounts/" + viper.GetString("CF_ID") + "/ai/run/@cf/meta/llama-3.3-70b-instruct-fp8-fast"
	authToken := "Bearer " + viper.GetString("CF_TOKEN")
	requestBody := ChatRequest{Messages: chatHistory}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("Error unmarshalling response: %w", err)
	}

	if !chatResponse.Success {
		return "", fmt.Errorf("API Error: %v", chatResponse.Errors)
	}

	return chatResponse.Result.Response, nil
}

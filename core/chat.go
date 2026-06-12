package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/model"
	"github.com/smtdfc/nagare/core/tool"
)

func chat() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	type GetWeatherArgs struct {
		Location string `json:"location_1" jsonschema:"description=Tên thành phố hoặc địa danh để tra cứu thời tiết, ví dụ: Hà Nội"`
	}

	tools := tool.ListTool{
		tool.DeclareTool("get_weather", "Lấy thông tin thời tiết tại địa điểm cho trước", func(args GetWeatherArgs) (any, error) {
			return "0 độ C, ban đêm trời nắng", nil
		}),
	}

	chatModel := model.NewOpenAICompatibleClient(&model.ChatModelConfig{
		BaseURL: "https://api.groq.com/openai/v1",
		APIKey:  os.Getenv("TOKEN"),
		Model:   "openai/gpt-oss-120b",
	})

	agentPool := agent.NewAgentPool(5, chatModel, tools)
	scanner := bufio.NewScanner(os.Stdin)
	sessionMgr := agent.NewSessionManager()
	sessionID := "1"

	for {
		ctx := context.Background()

		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("Nagare đi ngủ đây, bái bai! 🤡")
			break
		}

		a := agentPool.GetOrNew(chatModel)
		a.History = sessionMgr.GetHistory(sessionID)

		ch := a.Invoke(ctx, input)
		for msg := range ch {
			switch m := msg.(type) {
			case *messages.ToolCallMessage:
				fmt.Println(m.FunctionName)
				fmt.Println(m.Args)
			case *messages.ToolResultMessage:
				fmt.Println(m.CallID)
				fmt.Println(m.Error)
			case *messages.TextMessage:
				fmt.Print(m.Content)
			}
		}

		sessionMgr.SaveHistory(sessionID, a.History)
		fmt.Printf("\nHistory length: %d\n", len(a.History))

		for _, m := range a.History {
			fmt.Println(m)
		}

		agentPool.Put(a)
		fmt.Println()
	}

}

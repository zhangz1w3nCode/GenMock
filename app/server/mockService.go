package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAI API 密钥
const openAIKey = "sk-nA1PDr7FRRlxP6IL11FeDcCcA1344a2c91029d0d86F6Cd01"

// MockHandler 公共的 MockHandler 函数
func MockHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RequestType         string                 `json:"RequestType"`
		RequestMethodName   string                 `json:"RequestMethodName"`
		RequestMethodParams map[string]interface{} `json:"RequestMethodParams"`
	}

	// 解析请求体
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 生成模拟数据
	mockData := generateMockData(requestBody.RequestMethodName, requestBody.RequestType, requestBody.RequestMethodParams)

	// 返回模拟数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockData)
}

// 模拟数据生成函数
func generateMockData(typeStr, method string, params map[string]interface{}) interface{} {
	// 生成通用的 mock 数据
	mockData := make(map[string]interface{})

	// 生成请求内容
	requestContent := map[string]interface{}{
		"type":   typeStr,
		"method": method,
		"params": params,
	}

	// 调用 ChatGPT 生成 mock 数据
	generatedData, err := callChatGPT(requestContent)
	if err != nil {
		mockData["error"] = "Failed to generate mock data"
		return mockData
	}

	mockData["data"] = generatedData
	return mockData
}

// 调用 ChatGPT 生成数据
func callChatGPT(requestContent map[string]interface{}) (interface{}, error) {
	client := openai.NewClient(openAIKey) // 创建 OpenAI 客户端

	// 创建请求
	response, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: generatePrompt(requestContent),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	fmt.Println(response.Choices[0].Message.Content)

	var generatedData interface{}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &generatedData); err != nil {
		return nil, err
	}

	return generatedData, nil
}

// 生成请求的提示
func generatePrompt(requestContent map[string]interface{}) string {
	return "请你作为高级后端开发工程师 根据以下请求生成符合请求方法名称同时有意义的mock数据:\n" +
		"请求类型: " + requestContent["type"].(string) + "\n" +
		"请求方法: " + requestContent["method"].(string) + "\n" +
		"请求参数: " + fmt.Sprintf("%v", requestContent["params"]) + "\n" +
		"请返回一个 JSON 格式的mock数据"
}

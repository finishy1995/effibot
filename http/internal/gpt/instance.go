package gpt

import (
	"context"
	"fmt"
	"github.com/finishy1995/effibot/http/internal/config"
	"github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
)

// Handler is a gpt handler.
type Handler struct {
	client   *openai.Client
	maxToken int
}

// Message is a message for sending to OpenAI.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewHandler returns a new Handler.
func NewHandler(config *config.GPTConfig) *Handler {
	if config == nil {
		return nil
	}

	return &Handler{
		client:   newClient(config),
		maxToken: config.MaxToken,
	}
}

func newClient(config *config.GPTConfig) *openai.Client {
	// mock mode
	// mock 模式
	if config.Token == "" {
		return nil
	}

	c := openai.DefaultConfig(config.Token)

	var transport *http.Transport

	if config.TransportUrl != "" {
		// create a HTTP Transport object and set the proxy server
		// 创建一个 HTTP Transport 对象，并设置代理服务器
		proxyUrl, err := url.Parse(config.TransportUrl)
		if err != nil {
			panic(err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	// create a HTTP client and set the Transport object to its Transport field
	// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
	c.HTTPClient.Timeout = config.Timeout
	if transport != nil {
		c.HTTPClient.Transport = transport
	}

	return openai.NewClientWithConfig(c)
}

// CreateChatCompletion creates a chat completion.
func (h *Handler) CreateChatCompletion(ctx context.Context, messages []*Message) ([]string, int, error) {
	if messages == nil {
		return nil, 0, fmt.Errorf("messages cannot be nil")
	}
	length := len(messages)
	result := make([]string, 0, 1)
	if length == 0 {
		return result, 0, nil
	}

	// mock mode
	// mock 模式
	if h.client == nil {
		mockResult := "Mock mode detected, please set OpenAI token in config file!  \nMock 模式检测到，请在配置文件中设置 OpenAI token！  \n"
		mockResult += messages[length-1].Content

		return append(result, mockResult), len(mockResult) * 2, nil
	}

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo, // TODO: config
		Messages:    make([]openai.ChatCompletionMessage, length, length),
		MaxTokens:   h.maxToken,
		Temperature: 0.1, // TODO: config
	}

	for i, item := range messages {
		req.Messages[i] = openai.ChatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		}
	}

	resp, err := h.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("CreateChatCompletion failed, error: %s", err.Error())
	}

	for _, choice := range resp.Choices {
		result = append(result, choice.Message.Content)
		logx.WithContext(ctx).Infof("using token \tQuestion: %d Answer: %d Total: %d", resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
		logx.WithContext(ctx).Debugf("receive OpenAI response:\t %+v", resp)
	}
	return result, resp.Usage.TotalTokens, nil
}

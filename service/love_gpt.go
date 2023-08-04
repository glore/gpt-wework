package service

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"gpt-wework/result"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	gogpt "github.com/sashabaranov/go-gpt3"
)

// 聊天缓存
var loveGPTChatCache = cache.New(5*time.Minute, 5*time.Minute)

type LoveGPT struct {
	client *gogpt.Client
	ctx    context.Context
	userId string
}

// LoveGPTChat 小恋
func LoveGPTChat(c *gin.Context) {
	question := c.PostForm("question")
	userId := c.PostForm("userId")
	if len(question) == 0 {
		result.Fail(c, result.ResponseJson{Msg: "err:question"})
		return
	}
	if len(userId) == 0 {
		result.Fail(c, result.ResponseJson{Msg: "err:userId"})
		return
	}
	ret, err := LoveGPTService(question, userId, weworkConversationSize)
	if err != nil {
		result.Fail(c, result.ResponseJson{Msg: err.Error()})
		return
	}
	result.Success(c, result.ResponseJson{Data: ret})
}

func LoveGPTService(question, userId string, size int) (string, error) {
	var messages []gogpt.ChatCompletionMessage
	key := fmt.Sprintf("cache:love_gpt:chat:%s", userId)
	data, found := loveGPTChatCache.Get(key)
	if found {
		messages = data.([]gogpt.ChatCompletionMessage)
	}
	messages = append(messages, gogpt.ChatCompletionMessage{
		Role:    "system",
		Content: question,
	})
	fmt.Println("用户id:"+userId, messages)
	pivot := size
	if pivot > len(messages) {
		pivot = len(messages)
	}
	messages = messages[len(messages)-pivot:]
	loveGPTChatCache.Set(key, messages, 12*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		cancel()
	}()
	chat := &LoveGPT{
		client: gogpt.NewClient(openAiKey),
		ctx:    ctx,
		userId: userId,
	}

	defer chat.ctx.Done()
	answer, err := chat.Request(messages)
	if err != nil {
		fmt.Print(err.Error())
	}
	return answer, err
}

func (c *LoveGPT) Request(messages []gogpt.ChatCompletionMessage) (answer string, err error) {
	var msg = gogpt.ChatCompletionMessage{}
	msg.Role = "system"
	req := gogpt.ChatCompletionRequest{
		Model:    gogpt.GPT3Dot5Turbo,
		Messages: messages,
	}
	resp, err := c.client.CreateChatCompletion(c.ctx, req)
	if err != nil {
		return "", err
	}
	answer = resp.Choices[0].Message.Content
	for len(answer) > 0 {
		if answer[0] == '\n' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return answer, err
}

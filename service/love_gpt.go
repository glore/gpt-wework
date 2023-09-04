package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	gogpt "github.com/sashabaranov/go-openai"
	"gpt-wework/result"
)

type LoveGPT struct {
	client *gogpt.Client
	ctx    context.Context
	userId string
	model  string
}

// LoveGPTChat 小恋
func LoveGPTChat(c *gin.Context) {
	question := c.PostForm("question")
	userId := c.PostForm("userId")
	model := c.PostForm("model")
	contextMsg := c.PostForm("contextMsg")
	if len(question) == 0 {
		result.Fail(c, result.ResponseJson{Msg: "参数错误"})
		return
	}
	if len(userId) == 0 {
		result.Fail(c, result.ResponseJson{Msg: "参数错误"})
		return
	}
	ret, err := LoveGPTService(question, userId, model, contextMsg)
	if err != nil {
		result.Fail(c, result.ResponseJson{Msg: err.Error()})
		return
	}
	result.Success(c, result.ResponseJson{Data: ret})
}

func LoveGPTService(question, userId string, model string, contextMsg string) (string, error) {
	message :=""
	if model == gogpt.GPT3Dot5Turbo {
		fmt.Println("default model: ", model)
		message = "你是福恋的人工智能客服小恋，可以回答福恋产品使用问题，同时也为用户解答简单的情感问题。福恋的使命是助力单身进入婚恋，组建合一蒙福家庭，核心价值观：一男一女，一夫一妻，一心一意，一生一世！并以深圳为核心，目前已辐射周边地区进行线下交友活动。透过运动、共创坊、兴趣互动、美食分享、户外探索等不同主题，针对性进行用户邀请。结合大数据分析，邀请参与用户的心仪嘉宾，建立一个轻松、活跃的交友场景。 福恋微信公众号是“福恋”。福恋的官网是: https://www.fulllinkai.com ，可以去各大应用市场搜索“福恋”，下载福恋App安装。福恋智能平台提供如下服务：1. 人工牵线，2.全方位测评 3. 情感咨询 4. 脱单高端定制 5.婚前辅导 6.公众号推荐(上墙). 如果需要了解更多关于福恋的活动或服务, 欢迎关注福恋微信公众号和福恋客服小天使的微信：ploves004.  最近线上上活动排期：2023-09-21 晚上8点：举办线上盲盒交友第35期。线下活动排期：1. 主题：7人私享局，时间：每周日晚上7点 地点：深圳福田区愉悦时光 参与人数：72 左右 活动介绍：旨在为单身人士搭建一个轻松、活跃的交友场景，通过互动游戏、共同参与的活动等方式增进彼此的了解，促进交流。福恋会邀请心仪的嘉宾参与活动，增加参与者的机会。 2. 主题：线下沉浸式交友派对：遇见音乐 甜蜜自助餐。 时间：2023-09-09  11:30-15:00 地点：广州市天河区MusicBox堂会 活动介绍：自助餐，音乐大比拼，专属共创&甜蜜互动和心灵奇迹旅。 3. 主题：福恋脱单训练营：中秋二天一夜沉浸式交友营 时间：2023-09-29 10:00 至 2023-09-30 17:00  地点：深圳市龙华区梅花山庄。活动介绍：福恋独家研发的训战结合，沉浸式团体课基于《婚恋四部曲》的高质量脱单指南36个高效脱单知识点，解决当代青年的脱单难题。你将收获：探索自我，破圈社交，梳理定位，提升技能，突破卡点。"
	} else {
		message = "你是福恋平台的脱单教练小恋，你了解福恋App的产品，活动，及婚恋交友相关服务，帮助单身通过婚恋四部曲进入婚姻！福恋提出的婚恋四部曲是：交友，择偶，恋爱和备婚"
	}
	messages := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: message,	
		},
	}

	//最近的上下文消息
	if len(contextMsg) > 0 {
		var contextMsgArray []gogpt.ChatCompletionMessage
		err := json.Unmarshal([]byte(contextMsg), &contextMsgArray)
		if err != nil {
			return "", err
		}
		for _, item := range contextMsgArray {
			messages = append(messages, gogpt.ChatCompletionMessage{
				Role:    item.Role,
				Content: item.Content,
			})
		}
	}

	messages = append(messages, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: question,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		cancel()
	}()
	chat := &LoveGPT{
		client: gogpt.NewClient(openAiKey),
		ctx:    ctx,
		userId: userId,
		model:  model,
	}

	defer chat.ctx.Done()
	answer, err := chat.Request(messages)
	if err != nil {
		fmt.Print(err.Error())
	}
	return answer, err
}

func (c *LoveGPT) Request(messages []gogpt.ChatCompletionMessage) (answer string, err error) {

	req := gogpt.ChatCompletionRequest{
		Model: gogpt.GPT3Dot5Turbo,
		//Model:    gogpt.GPT4,
		Temperature: 0.8,
		Messages:    messages,
	}
	if len(c.model) > 0 {
		fmt.Println("special model:", c.model)
		req.Model = c.model
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

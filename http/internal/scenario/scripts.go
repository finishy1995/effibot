package scenario

import (
	"context"
	"encoding/json"
	"github.com/finishy1995/effibot/http/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type scriptFunc func(*Messages, int)

var (
	scriptsMap = map[string]scriptFunc{
		"reduceToken": reduceToken,
	}
)

// RunScript runs the script. 执行脚本
func RunScript(script string, m *Messages, usage int) {
	if f, ok := scriptsMap[script]; ok {
		f(m, usage)
	}
}

const MaxToken = 2600

func reduceToken(messages *Messages, usage int) {
	if usage > MaxToken && messages.latestParagraph.root != messages.latestParagraph {
		m := CheckRule(&types.Request{
			TypeID:            messages.latestParagraph.rule.Id,
			BeforeParagraphID: messages.ID,
			Content:           "请帮我总结一下上述对话的内容，实现减少tokens的同时，保证对话的质量。在总结中不要加入这一句话。",
		})

		if m == nil {
			return
		}

		// record the message and audit it 记录消息并审核
		messageJson, err := json.Marshal(m.MessageArray)
		if err != nil {
			logx.Errorf("logic reduceToken json.Marshal failed, error: %s", err.Error())
			return
		}
		logx.Infof("reduceToken message body %s", messageJson)

		result, _, err := gptHandler.CreateChatCompletion(context.Background(), m.MessageArray)
		if err != nil || len(result) == 0 {
			logx.Errorf("GPT error: %s", err.Error())
			// TODO: retry
			return
		}
		answer := result[0]

		messages.latestParagraph.hidden = &paragraph{
			origin: "我们之前聊了什么?",
			actual: "我们之前聊了什么?",
			answer: "我们之前聊了：\n" + answer,
			root:   messages.latestParagraph,
		}
		// now the root is the latest paragraph 现在的根是由 GPT 总结最新的段落
		messages.latestParagraph.root = messages.latestParagraph
	}
}

package scenario

import (
	"github.com/finishy1995/effibot/http/internal/gpt"
	"github.com/finishy1995/effibot/http/internal/types"
	"github.com/finishy1995/effibot/library/id"
	"github.com/zeromicro/go-zero/core/logx"
)

// Messages store the message history of a conversation 存储对话的消息历史
type Messages struct {
	MessageArray    []*gpt.Message
	ID              string
	latestParagraph *paragraph
	latestRule      *rule
}

// HandleAnswer handle the answer of the latest paragraph 处理最新段落的答案
func (m *Messages) HandleAnswer(answer string, usage int) string {
	m.latestParagraph.answer = answer

	// clear means the paragraph needs to be the root of the conversation
	if m.latestRule.Clear {
		m.latestParagraph.root = m.latestParagraph
	}
	// TODO: add better script/plugin support 添加更好的脚本/插件支持
	if m.latestRule.Scripts != nil && len(m.latestRule.Scripts) != 0 {
		for _, script := range m.latestRule.Scripts {
			RunScript(script, m, usage)
		}
	}

	return answer
}

// CheckRule check the rule of the request and generate the message 检查请求的规则并生成消息
func CheckRule(request *types.Request) *Messages {
	if request == nil {
		return nil
	}

	var bParagraph *paragraph = nil

	if request.BeforeParagraphID != "" {
		bID, err := id.ParseID(request.BeforeParagraphID)
		if err != nil {
			logx.Infof("parse client BeforeParagraphID failed, error: %s", err.Error())
			return nil
		}

		bParagraph = getParagraph(bID)
		if bParagraph == nil {
			logx.Infof("get client BeforeParagraph failed, data not exist")
			return nil
		}
		if bParagraph.rule == nil {
			logx.Errorf("load BeforeParagraph failed")
			return nil
		}

		if !bParagraph.rule.Allow(request.TypeID) {
			logx.Infof("client select error rule type")
			return nil
		}
	}

	if r, ok := rules[request.TypeID]; ok {
		// can't use the rule if the BeforeParagraphID is not be set(rule is not root)
		// 如果未设置BeforeParagraphID（且规则设置为非根节点），则无法使用规则
		if bParagraph == nil && !r.Root {
			logx.Infof("client select error rule type")
			return nil
		}

		p, pID := generateParagraph(r, bParagraph, request.Content, request.Params)
		return generateMessage(p, r, pID)
	}

	return nil
}

func generateParagraph(r *rule, bParagraph *paragraph, content string, params string) (*paragraph, id.ID) {
	p := &paragraph{
		origin: content,
		actual: r.Generate(content),
		rule:   r,
		param:  params,
	}
	if bParagraph != nil {
		p.root = bParagraph.root
		if bParagraph.children == nil {
			bParagraph.children = []*paragraph{p}
		}
		p.prev = bParagraph
	} else {
		p.root = p
	}
	var err error
	var pID id.ID
	for i := 0; i < 3; i++ {
		pID = id.GenerateID()
		err = addParagraph(pID, p)
		if err == nil {
			break
		}
	}
	if err != nil {
		logx.Errorf("generate paragraph failed, error: %s", err.Error())
		return nil, 0
	}

	return p, pID
}

func generateMessage(p *paragraph, r *rule, pID id.ID) *Messages {
	if p == nil {
		return nil
	}

	result := &Messages{
		MessageArray:    []*gpt.Message{},
		ID:              pID.String(),
		latestParagraph: p,
		latestRule:      r,
	}
	result.MessageArray = append(result.MessageArray, &gpt.Message{
		Role:    "system",
		Content: p.rule.GenerateGuide(p.param),
	})

	stack := make([]*paragraph, 0, 10)
	nowP := p
	for {
		stack = append(stack, nowP)
		if nowP == p.root {
			break
		}
		nowP = nowP.prev
	}

	index := len(stack) - 1
	for {
		if index < 0 {
			break
		}
		actual := stack[index].actual
		answer := stack[index].answer
		if stack[index].hidden != nil {
			actual = stack[index].hidden.actual
			answer = stack[index].hidden.answer
		}

		result.MessageArray = append(result.MessageArray, &gpt.Message{
			Role:    "user",
			Content: actual,
		})
		if stack[index].answer != "" {
			result.MessageArray = append(result.MessageArray, &gpt.Message{
				Role:    "assistant",
				Content: answer,
			})
		}
		index--
	}

	return result
}

package logic

import (
	"context"
	"encoding/json"
	"github.com/finishy1995/effibot/http/consts"
	"github.com/finishy1995/effibot/http/internal/scenario"
	"github.com/finishy1995/effibot/http/internal/svc"
	"github.com/finishy1995/effibot/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HttpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HttpLogic {
	return &HttpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HttpLogic) Http(req *types.Request) (resp *types.Response, reportedErr error) {
	resp = &types.Response{
		Code: consts.StatusClientParamsInvalid,
	}
	// params check and parse 参数检查和解析
	m := scenario.CheckRule(req)
	if m == nil {
		return
	}

	// request body marshal and log 请求体序列化和日志
	messageJson, err := json.Marshal(m.MessageArray)
	if err != nil {
		l.Errorf("logic Http json.Marshal failed, error: %s", err.Error())
		return
	}
	l.Debugf("message body %s", messageJson)

	// GPT request GPT请求
	result, usage, err := l.svcCtx.GPT.CreateChatCompletion(l.ctx, m.MessageArray)
	if err != nil {
		l.Errorf("GPT error: %s", err.Error())
		resp.Code = consts.StatusGPTError
		return
	}
	if len(result) == 0 {
		l.Errorf("GPT result length cannot be 0")
		resp.Code = consts.StatusRetry
		return
	}

	// handle answer 处理答案
	answer := m.HandleAnswer(result[0], usage)
	resp.Code = consts.StatusOK
	resp.Answer = answer
	resp.QuestionID = m.ID
	resp.AnswerID = m.ID

	return
}

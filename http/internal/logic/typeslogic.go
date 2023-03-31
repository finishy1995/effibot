package logic

import (
	"context"
	"github.com/finishy1995/effibot/http/internal/scenario"

	"github.com/finishy1995/effibot/http/internal/svc"
	"github.com/finishy1995/effibot/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TypesLogic {
	return &TypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TypesLogic) Types(_ *types.TypesRequest) (resp *types.TypesResponse, err error) {
	resp = &types.TypesResponse{
		Types: scenario.GetAllRules(),
	}
	return
}

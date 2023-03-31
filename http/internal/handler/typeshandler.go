package handler

import (
	"net/http"

	"github.com/finishy1995/effibot/http/internal/logic"
	"github.com/finishy1995/effibot/http/internal/svc"
	"github.com/finishy1995/effibot/http/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TypesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(allowOrigin, allOrigins)
		w.Header().Set(allowMethods, methods)
		w.Header().Set(allowHeaders, headers)
		
		var req types.TypesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTypesLogic(r.Context(), svcCtx)
		resp, err := l.Types(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

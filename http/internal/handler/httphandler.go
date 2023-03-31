package handler

import (
	"net/http"

	"github.com/finishy1995/effibot/http/internal/logic"
	"github.com/finishy1995/effibot/http/internal/svc"
	"github.com/finishy1995/effibot/http/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	allowOrigin  = "Access-Control-Allow-Origin"
	allOrigins   = "*"
	allowMethods = "Access-Control-Allow-Methods"
	allowHeaders = "Access-Control-Allow-Headers"
	headers      = "Content-Type, Content-Length, Origin"
	methods      = "GET, HEAD, POST, PATCH, PUT, DELETE"
)

func HttpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(allowOrigin, allOrigins)
		w.Header().Set(allowMethods, methods)
		w.Header().Set(allowHeaders, headers)

		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewHttpLogic(r.Context(), svcCtx)
		resp, err := l.Http(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

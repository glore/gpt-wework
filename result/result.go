package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ResponseJson struct {
	Status int         `json:"-"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func (r ResponseJson) IsEmpty() bool {
	return reflect.DeepEqual(r, ResponseJson{})
}

func buildStatus(resp ResponseJson, defaultStatus int) int {
	if resp.Status == 0 {
		return defaultStatus
	}
	return resp.Status
}

func HttpResponse(ctx *gin.Context, status int, resp ResponseJson) {
	if resp.IsEmpty() {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.AbortWithStatusJSON(status, resp)
}

func Success(ctx *gin.Context, resp ResponseJson) {
	if reflect.ValueOf(resp.Msg).IsZero() {
		resp.Msg = "success"
	}
	if reflect.ValueOf(resp.Code).IsZero() {
		resp.Code = 0
	}
	HttpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}

func Fail(ctx *gin.Context, resp ResponseJson) {
	if reflect.ValueOf(resp.Msg).IsZero() {
		resp.Msg = "fail"
	}
	if reflect.ValueOf(resp.Code).IsZero() {
		resp.Code = 1
	}
	HttpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}

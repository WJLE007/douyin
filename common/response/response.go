package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
	Data any    `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp any, err error) {
	var body Body
	if err != nil {
		body.Code = -1
		body.Msg = err.Error()
	} else {
		body.Msg = "Success"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

func Error(w http.ResponseWriter, err error) {
	var body Body
	body.Code = -1
	body.Msg = err.Error()
	httpx.OkJson(w, body)
}

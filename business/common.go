package business

import (
	"fmt"
	"net/http"

	"github.com/samuelsih/goth/model"
)

type CommonRequest struct {
	Session model.UserSession
}

type CommonResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"messsage"`
}

func (o *CommonResponse) SetError(code int, msg string) {
	o.Code = code
	o.Msg = msg
	if o.Code >= 500 {
		fmt.Printf("ERROR %d = %v\n", o.Code, o.Msg)
	}
}

func (o *CommonResponse) SetOK() {
	o.Code = http.StatusOK
	o.Msg = "ok"
}

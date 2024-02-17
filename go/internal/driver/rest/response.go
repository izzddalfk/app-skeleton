package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// ADJUST: feel free to adjust this struct following your api specification
type RespBody struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	Err        string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
	Timestamp  int64       `json:"ts"`
}

func (b *RespBody) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, b.StatusCode)
	b.Timestamp = time.Now().Unix()

	return nil
}

func NewSuccessResp(data interface{}) *RespBody {
	return &RespBody{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
	}
}

func NewErrorResp(err error) *RespBody {
	restErr := readErr(err)

	return &RespBody{
		OK:         false,
		StatusCode: restErr.StatusCode,
		Err:        restErr.Err,
		Message:    restErr.Message,
	}
}

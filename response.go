package ysx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Code int
	Msg  string
	Data json.RawMessage
}

type ErrResponse struct {
	Code int
	Msg  string
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

func (c *client) newResponse(r *http.Response) (*Response, error) {
	rsp := new(Response)
	if r.Body != nil {
		slurp, err := ioutil.ReadAll(io.Reader(r.Body))
		if err != nil {
			return nil, err
		}
		if len(slurp) > 0 {
			err := json.Unmarshal(slurp, rsp)
			if err != nil {
				return nil, err
			}
			if rsp.Code == 10403 {
				return nil, errors.New("token过期")
			}
			if rsp.Code != 200 {
				e := &ErrResponse{
					Code: rsp.Code,
					Msg:  rsp.Msg,
				}

				return nil, e
			}
		}
	}

	return rsp, nil
}

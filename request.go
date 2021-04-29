package ysx

import (
	`bytes`
	`encoding/json`
	`fmt`
	`io`
	`io/ioutil`
	`net/http`
	`strings`
)

type Request http.Request

func NewRequest(method, url string) (*Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return (*Request)(req), nil
}

func (r *Request) SetBody(body interface{}) error {
	switch b := body.(type) {
	case string:
		return r.setBodyString(b)
	default:
		return r.setBodyJson(b)
	}
}

func (r *Request) setBodyJson(body interface{}) error {
	data, err := json.Marshal(body)
	s := string(data)
	fmt.Println(s)
	if err != nil {
		return err
	}

	return r.setBodyReader(bytes.NewReader(data))
}

func (r *Request) setBodyString(body string) error {
	return r.setBodyReader(strings.NewReader(body))
}

func (r *Request) setBodyReader(body io.Reader) error {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *strings.Reader:
			r.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			r.ContentLength = int64(v.Len())
		}
	}

	return nil
}

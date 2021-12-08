package ysx

import (
	`encoding/json`
	`fmt`
	`net/http`
	`net/url`
	`time`
)

type token interface {
	GetToken()
}

type TokenRsp struct {
	AccessToken             string           `json:"accessToken"`
	RefreshToken            *string          `json:"refreshToken"`
	ExpiresTime             string           `json:"expiresTime"`
	UserId                  string           `json:"userId"`
	UserName                string           `json:"userName"`
	EnterpriseId            string           `json:"enterpriseId"`
	BusinessType            []int            `json:"businessType"`
	Roles                   []string         `json:"roles"`
	Authorities             map[string]*bool `json:"authorities"`
	DiscussionFlag          bool             `json:"discussionFlag"`
	WeakPwd                 bool             `json:"weakPwd"`
	SourceType              int              `json:"sourceType"`
	ConfidentialityFunction *bool            `json:"confidentialityFunction"`
	VoiceAiFlag             int              `json:"voiceAiFlag"`
	OpenStatus              *bool            `json:"openStatus"`
	IfFirstLogin            *bool            `json:"if_first_login"`
}

func (c *client) GetToken() {
	// 判断是否有其他goroutine在拿token
	// 如果有，还是阻塞在这
	do := len(c.ch) == 0

	c.ch <- struct{}{}
	defer func() {
		<-c.ch
	}()

	// 如果最开始有其他goroutine在拿token，则此处直接退出
	if !do {
		return
	}

	var (
		token string
		rsp   *TokenRsp
	)

	for i := 0; i < c.retryLimit+1; i++ {
		r, err := c.performRequest(PerformRequestOptions{
			Method: http.MethodGet,
			Path:   "/mixapi/token",
			Params: url.Values{
				"identity": []string{c.basicIdentity},
				"mobile":   []string{c.basicMobile},
				"key":      []string{c.basicKey},
			},
		})
		if err != nil {
			fmt.Println(fmt.Sprintf("---获取token失败：%s，重试第%d次", err, i))

			continue
		}

		rsp = new(TokenRsp)
		err = json.Unmarshal(r.Data, rsp)
		if err != nil {
			fmt.Println(fmt.Sprintf("---获取token失败：%s，重试第%d次", err, i))

			continue
		}

		if rsp.RefreshToken != nil {
			token = *rsp.RefreshToken
		} else {
			token = rsp.AccessToken
		}

		c.header["X-ACCESS-TOKEN"] = []string{token}

		fmt.Println("---获取token成功")

		return
	}

	panic("---获取token失败，程序退出")

	return
}

func refreshToken(c *client) {
	for {
		time.Sleep(c.tokenRefreshInterval)

		c.GetToken()
	}
}

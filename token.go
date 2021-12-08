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
	c.mu.Lock()
	defer c.mu.Unlock()

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

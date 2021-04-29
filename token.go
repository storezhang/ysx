package ysx

import (
	`encoding/json`
	`fmt`
	`net/http`
	`net/url`
	`time`
)

type token interface {
	GetToken() (*TokenRsp, error)
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

func (c *client) GetToken() (*TokenRsp, error) {
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
		return nil, err
	}

	rsp := new(TokenRsp)
	err = json.Unmarshal(r.Data, rsp)
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	if rsp.RefreshToken != nil {
		c.header["X-ACCESS-TOKEN"] = []string{*rsp.RefreshToken}
	} else {
		c.header["X-ACCESS-TOKEN"] = []string{rsp.AccessToken}
	}
	c.mu.RUnlock()

	return rsp, nil
}

func refreshToken(c *client) {
	for {
		time.Sleep(c.tokenRefreshInterval)

		_, err := c.GetToken()
		if err != nil {
			fmt.Println(err)
		}
	}
}

package ysx

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type getMeetingInfoService struct {
	c *client

	meetingId string
}

type MeetingInfoRsp struct {
	DeputyChairsControlSwitch bool        `json:"deputychairs_control_switch"`
	IfInviteHost              bool        `json:"if_invite_host"`
	IfOpenAiTrans             bool        `json:"if_open_aitrans"`
	IsDeputyChairs            bool        `json:"is_deputychairs"`
	MeetingContent            string      `json:"meeting_content"`
	MeetingCountDown          string      `json:"meeting_countDown"`
	MeetingStartTime          string      `json:"meeting_starttime"`
	MeetingEndTime            string      `json:"meeting_endtime"`
	MeetingFrom               int         `json:"meeting_from"`
	MeetingHostId             string      `json:"meeting_host_id"`
	MeetingId                 string      `json:"meeting_id"`
	MeetingIdReturn           string      `json:"meeting_id_return"`
	MeetingIfLive             bool        `json:"meeting_iflive"`
	MeetingIfLock             bool        `json:"meeting_iflock"`
	MeetingIfMute             bool        `json:"meeting_ifmute"`
	MeetingIfRecord           bool        `json:"meeting_ifrecord"`
	MeetingLength             int64       `json:"meeting_length"`
	MeetingMode               meetingMode `json:"meeting_mode"`
	MeetingOrderLength        int         `json:"meeting_order_length"`
	MeetingPassword           string      `json:"meeting_password"`
	MeetingStatus             string      `json:"meeting_status"`
	MeetingTheme              string      `json:"meeting_theme"`
	MeetingType               string      `json:"meeting_type"`
	MeetingVoiceMode          string      `json:"meeting_voicemode"`
	OperatorId                string      `json:"operator_id"`
	SupportChangeTokenOwner   bool        `json:"support_change_token_owner"`
	WebsocketUrl              string      `json:"websocket_url"`
	WxMiniSwitch              bool        `json:"wxmini_switch"`
}

func (c *client) GetMeetingInfo() *getMeetingInfoService {
	return &getMeetingInfoService{
		c: c,
	}
}

func (g *getMeetingInfoService) MeetingId(meetingId string) *getMeetingInfoService {
	g.meetingId = meetingId

	return g
}

func (g *getMeetingInfoService) Do() (*MeetingInfoRsp, error) {
	r, err := g.c.performRequest(PerformRequestOptions{
		Method: http.MethodGet,
		Path:   "/conference/meeting/getMeetingInfo",
		Params: url.Values{
			"meeting_id": []string{g.meetingId},
		},
		Body:    nil,
		Headers: nil,
	})
	if err != nil {
		return nil, err
	}

	rsp := new(MeetingInfoRsp)
	err = json.Unmarshal(r.Data, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

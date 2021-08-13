package ysx

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	MeetingOperateTypePublic meetingOperateType = 0
	MeetingOperateTypeSecret meetingOperateType = 1
	MeetingOperateTypeSoft   meetingOperateType = 2

	MeetingTypeAppointment meetingType = 0
	MeetingTypeInstance    meetingType = 1

	MeetingModeHighDefinition meetingMode = "HD"
	MeetingModeAudio          meetingMode = "VO"

	LiveTypePublic  liveType = "PUBLIC"
	LiveTypePrivate liveType = "PRIVATE"

	defaultMeetingLength = 45 * 60
)

type (
	meetingOperateType int8
	meetingType        int8
	meetingMode        string
	liveType           string
)

type meeting interface {
	CreateMeeting() *createMeetingService
	EndMeeting() *endMeetingService
	InviteAttendee() *inviteAttendeeService
	SetChairman() *setChairmanService
	GetMeetingDevice() *getMeetingDeviceService
	MultiScreenSetting() *multiScreenSettingService
	ReInvite() *reInviteService
	GetMeetingInfo() *getMeetingInfoService
}

type createMeetingService struct {
	c *client

	meetingOperateType meetingOperateType
	meetingTheme       string
	meetingAttendee    []*Attendee
	meetingContent     string
	meetingIfMute      bool
	meetingLength      int
	meetingMode        meetingMode
	meetingStartTime   string
	meetingType        meetingType
	extendData         interface{}
	ifInviteHost       bool
	body               *CreateMeetingReq
	bodyString         string
}

type PublicMeetingExtendData struct {
	MeetingIfLive bool     `json:"meeting_if_live"`
	LiveAccessKey string   `json:"live_accesskey,omitempty"`
	LiveType      liveType `json:"live_type,omitempty"`
}

type SoftMeetingExtendData struct {
	OpenHostVideo bool `json:"open_host_video"`
}

type Attendee struct {
	Mobile  string `json:"mobile"`
	IsVolte bool   `json:"is_volte"`
}

type CreateMeetingReq struct {
	MeetingOperateType meetingOperateType `json:"meeting_operate_type"`
	MeetingTheme       string             `json:"meeting_theme"`
	MeetingAttendee    []*Attendee        `json:"meeting_attendee,omitempty"`
	MeetingContent     *string            `json:"meeting_content,omitempty"`
	MeetingIfMute      *bool              `json:"meeting_ifmute,omitempty"`
	MeetingLength      int                `json:"meeting_length"`
	MeetingMode        meetingMode        `json:"meeting_mode"`
	MeetingStartTime   *string            `json:"meeting_starttime,omitempty"`
	MeetingType        meetingType        `json:"meeting_type"`
	ExtendData         interface{}        `json:"extend_data,omitempty"`
	IfInviteHost       *bool              `json:"if_invite_host,omitempty"`
}

func (c *client) CreateMeeting() *createMeetingService {
	return &createMeetingService{
		c:                  c,
		meetingOperateType: MeetingOperateTypePublic,
		meetingTheme:       "",
		meetingAttendee:    nil,
		meetingContent:     "",
		meetingIfMute:      false,
		meetingLength:      defaultMeetingLength,
		meetingMode:        MeetingModeHighDefinition,
		meetingStartTime:   "",
		meetingType:        MeetingTypeInstance,
		extendData:         nil,
		ifInviteHost:       false,
	}
}

func (c *createMeetingService) OperateType(typ meetingOperateType) *createMeetingService {
	c.meetingOperateType = typ

	return c
}

func (c *createMeetingService) Theme(theme string) *createMeetingService {
	c.meetingTheme = theme

	return c
}

func (c *createMeetingService) Attendee(attendees ...*Attendee) *createMeetingService {
	c.meetingAttendee = append(c.meetingAttendee, attendees...)

	return c
}

func (c *createMeetingService) Content(content string) *createMeetingService {
	c.meetingContent = content

	return c
}

func (c *createMeetingService) Mute(mute bool) *createMeetingService {
	c.meetingIfMute = mute

	return c
}

func (c *createMeetingService) Last(second int) *createMeetingService {
	if second > 1 {
		c.meetingLength = second
	}

	return c
}

func (c *createMeetingService) Mode(mode meetingMode) *createMeetingService {
	c.meetingMode = mode

	return c
}

func (c *createMeetingService) StartAt(t interface{}) *createMeetingService {
	switch t := t.(type) {
	case time.Time:
		c.meetingStartTime = t.Format("2006-01-02 15:04:05")
	case string:
		c.meetingStartTime = t
	}

	return c
}

func (c *createMeetingService) Type(typ meetingType) *createMeetingService {
	c.meetingType = typ

	return c
}

func (c *createMeetingService) ExtendData(data interface{}) *createMeetingService {
	switch data.(type) {
	case *PublicMeetingExtendData, *SoftMeetingExtendData:
		c.extendData = data
	}

	return c
}

func (c *createMeetingService) InviteHost(invite bool) *createMeetingService {
	c.ifInviteHost = invite

	return c
}

func (c *createMeetingService) Body(body *CreateMeetingReq) *createMeetingService {
	c.body = body

	return c
}

func (c *createMeetingService) BodyString(body string) *createMeetingService {
	c.bodyString = body

	return c
}

func (c *createMeetingService) buildBody() (interface{}, error) {
	if c.body != nil {
		return c.body, nil
	}
	if c.bodyString != "" {
		return c.bodyString, nil
	}
	if len(c.meetingAttendee) > 50 {
		c.meetingAttendee = c.meetingAttendee[:50]
	}

	req := &CreateMeetingReq{
		MeetingOperateType: c.meetingOperateType,
		MeetingTheme:       c.meetingTheme,
		MeetingAttendee:    c.meetingAttendee,
		MeetingLength:      c.meetingLength,
		MeetingMode:        c.meetingMode,
		MeetingType:        c.meetingType,
		ExtendData:         c.extendData,
	}

	if c.meetingType == MeetingTypeAppointment {
		req.MeetingStartTime = &c.meetingStartTime
	}
	if c.meetingOperateType != MeetingOperateTypeSoft {
		req.MeetingContent = &c.meetingContent
		req.MeetingIfMute = &c.meetingIfMute
	}
	switch req.MeetingOperateType {
	case MeetingOperateTypePublic:
		req.MeetingMode = c.meetingMode
		req.IfInviteHost = &c.ifInviteHost
		req.ExtendData = c.extendData
	case MeetingOperateTypeSecret:
		req.MeetingMode = MeetingModeHighDefinition
	case MeetingOperateTypeSoft:
		req.MeetingMode = MeetingModeHighDefinition
		req.ExtendData = c.extendData
	}

	return req, nil
}

func (c *createMeetingService) Do() (string, error) {
	req, err := c.buildBody()
	if err != nil {
		return "", err
	}

	rsp, err := c.c.performRequest(PerformRequestOptions{
		Method:  http.MethodPost,
		Path:    "/mixapi/createMeeting",
		Params:  nil,
		Body:    req,
		Headers: nil,
	})
	if err != nil {
		return "", err
	}

	var s string
	err = json.Unmarshal(rsp.Data, &s)
	if err != nil {
		return "", err
	}

	return s, nil
}

package ysx

import (
	`encoding/json`
	`net/http`
)

type multiScreenSettingService struct {
	c *client

	imageType                   string
	meetingId                   string
	meetingScreenSettingRestore string
	meetingScreenType           string
	speakerCallId               string
	subscribers                 []*Subscriber
	switchTime                  int
	chosenDevices               []*ChosenDevice
	body                        *MultiScreenSettingReq
	bodyString                  string
}

type MultiScreenSettingReq struct {
	ImageType                   string        `json:"image_type"`
	MeetingId                   string        `json:"meeting_id"`
	MeetingScreenSettingRestore string        `json:"meeting_screensettingrestore"`
	MeetingScreenType           string        `json:"meeting_screentype"`
	SpeakerCallId               string        `json:"speaker_call_id,omitempty"`
	Subscribers                 []*Subscriber `json:"subscribers"`
	SwitchTime                  int           `json:"switch_time"`
	IfUseTemplate               bool          `json:"if_use_template"`
}

type meetingScreenSettingRestore struct {
	ImageType         string          `json:"image_type"`
	MeetingScreenType string          `json:"meeting_screentype"`
	ChosenDevices     []*ChosenDevice `json:"choosenDevices"`
}

type ChosenDevice struct {
	DeviceId         string  `json:"device_id"`
	DeviceName       string  `json:"device_name"`
	MeetingCallId    string  `json:"meeting_call_id"`
	DeviceDepartName string  `json:"device_depart_name"`
	SipNo            *string `json:"sip_no"`
	EnterpriseUserId string  `json:"enterprise_userid"`
	EnterpriseName   string  `json:"enterprise_name"`
	IsChosen         bool    `json:"ischoosen"`
	Index            int     `json:"index"`
}

type Subscriber struct {
	Subscriber    string   `json:"subscriber"`
	AssistStream  bool     `json:"assist_stream"`
	DeviceId      string   `json:"device_id"`
	DeviceName    string   `json:"device_name"`
	DisplayCallId []string `json:"display_callid"`
	Index         int      `json:"index"`
	MainPicture   bool     `json:"main_picture"`
}

func (c *client) MultiScreenSetting() *multiScreenSettingService {
	return &multiScreenSettingService{
		c: c,

		imageType:         "Single",
		meetingScreenType: "static",
	}
}

func (m *multiScreenSettingService) ImageType(typ string) *multiScreenSettingService {
	m.imageType = typ

	return m
}

func (m *multiScreenSettingService) MeetingId(id string) *multiScreenSettingService {
	m.meetingId = id

	return m
}

func (m *multiScreenSettingService) ScreenType(typ string) *multiScreenSettingService {
	m.meetingScreenType = typ

	return m
}

func (m *multiScreenSettingService) SpeakerCallId(id string) *multiScreenSettingService {
	m.speakerCallId = id

	return m
}

func (m *multiScreenSettingService) Subscribers(subscribers ...*Subscriber) *multiScreenSettingService {
	m.subscribers = append(m.subscribers, subscribers...)

	return m
}

func (m *multiScreenSettingService) SwitchTime(t int) *multiScreenSettingService {
	m.switchTime = t

	return m
}

func (m *multiScreenSettingService) ChosenDevices(devices ...*ChosenDevice) *multiScreenSettingService {
	m.chosenDevices = append(m.chosenDevices, devices...)

	return m
}

func (m *multiScreenSettingService) Body(body *MultiScreenSettingReq) *multiScreenSettingService {
	m.body = body

	return m
}

func (m *multiScreenSettingService) BodyString(body string) *multiScreenSettingService {
	m.bodyString = body

	return m
}

func (m *multiScreenSettingService) buildBody() (interface{}, error) {
	if m.body != nil {
		return m.body, nil
	}
	if m.bodyString != "" {
		return m.bodyString, nil
	}

	for i := range m.chosenDevices {
		m.chosenDevices[i].Index = i + 1
	}

	restore := &meetingScreenSettingRestore{
		ImageType:         m.imageType,
		MeetingScreenType: m.meetingScreenType,
		ChosenDevices:     m.chosenDevices,
	}
	b, err := json.Marshal(restore)
	if err != nil {
		return nil, err
	}

	req := &MultiScreenSettingReq{
		ImageType:                   m.imageType,
		MeetingId:                   m.meetingId,
		MeetingScreenSettingRestore: string(b),
		MeetingScreenType:           m.meetingScreenType,
		SpeakerCallId:               m.speakerCallId,
		Subscribers:                 m.subscribers,
		SwitchTime:                  m.switchTime,
		IfUseTemplate:               true,
	}

	return req, nil
}

func (m *multiScreenSettingService) Do() error {
	req, err := m.buildBody()
	if err != nil {
		return err
	}

	_, err = m.c.performRequest(PerformRequestOptions{
		Method:  http.MethodPost,
		Path:    "/conference/meeting/multiScreenSettings",
		Params:  nil,
		Body:    req,
		Headers: nil,
	})

	return err
}

package ysx

import (
	`encoding/json`
	`net/http`
	`net/url`
	`strconv`
)

type getMeetingDeviceService struct {
	c *client

	meetingOperateType meetingOperateType
	meetingId          string
	searchKey          string
	view               string
}

type MeetingDeviceRsp struct {
	UserList       []*JoinMember `json:"user_list"`
	UserTotal      int           `json:"user_total"`
	EnterUserTotal int           `json:"enter_user_total"`
}

type JoinMember struct {
	SerialNo                string  `json:"serialno"`
	DeviceId                string  `json:"device_id"`
	Status                  string  `json:"status"`
	IfRollCall              bool    `json:"if_rollcall"`
	IfAttend                bool    `json:"if_attend"`
	IfVoice                 bool    `json:"if_voice"`
	MeetingCallId           string  `json:"meeting_call_id"`
	DeviceAssistStatus      bool    `json:"device_assist_status"`
	MediaType               string  `json:"media_type"`
	DeviceAuthority         string  `json:"device_authority"`
	DeviceDepartName        string  `json:"device_depart_name"`
	DepartmentId            string  `json:"department_id"`
	EnterpriseName          string  `json:"enterprise_name"`
	EnterpriseId            string  `json:"enterprise_id"`
	IsVip                   bool    `json:"is_vip"`
	IfVolteInvite           bool    `json:"if_volte_invite"`
	DeviceName              string  `json:"device_name"`
	DeviceType              string  `json:"device_type"`
	SipNo                   *string `json:"sip_no"`
	EnterpriseUserId        string  `json:"enterprise_userid"`
	InvitationFailedTimes   int     `json:"invitation_failed_times"`
	IfRadioSource           bool    `json:"is_radio_source"`
	UserDefinedOrder        int     `json:"user_defined_order"`
	ConfidentialityFunction bool    `json:"confidentiality_function"`
}

func (c *client) GetMeetingDevice() *getMeetingDeviceService {
	return &getMeetingDeviceService{
		c: c,

		meetingOperateType: MeetingOperateTypePublic,
		view:               "2",
	}
}

func (g *getMeetingDeviceService) OperateType(typ meetingOperateType) *getMeetingDeviceService {
	g.meetingOperateType = typ

	return g
}

func (g *getMeetingDeviceService) Id(id string) *getMeetingDeviceService {
	g.meetingId = id

	return g
}

func (g *getMeetingDeviceService) Search(key string) *getMeetingDeviceService {
	g.searchKey = key

	return g
}

func (g *getMeetingDeviceService) Do() (*MeetingDeviceRsp, error) {
	r, err := g.c.performRequest(PerformRequestOptions{
		Method: http.MethodGet,
		Path:   "/mixapi/getMeetingDevices",
		Params: url.Values{
			"view":                 []string{g.view},
			"search_key":           []string{g.searchKey},
			"meeting_id":           []string{g.meetingId},
			"meeting_operate_type": []string{strconv.Itoa(int(g.meetingOperateType))},
		},
		Body:    nil,
		Headers: nil,
	})
	if err != nil {
		return nil, err
	}

	rsp := new(MeetingDeviceRsp)
	err = json.Unmarshal(r.Data, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

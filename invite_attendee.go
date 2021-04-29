package ysx

import `net/http`

type inviteAttendeeService struct {
	c *client

	meetingOperateType meetingOperateType
	meetingId          string
	attendee           []*Attendee
	body               *InviteAttendeeReq
	bodyString         string
}

type InviteAttendeeReq struct {
	MeetingOperateType meetingOperateType `json:"meeting_operate_type"`
	MeetingId          string             `json:"meeting_id"`
	JoinUser           []*Attendee        `json:"join_user"`
}

func (c *client) InviteAttendee() *inviteAttendeeService {
	return &inviteAttendeeService{
		c: c,

		meetingOperateType: MeetingOperateTypePublic,
	}
}

func (i *inviteAttendeeService) OperateType(typ meetingOperateType) *inviteAttendeeService {
	i.meetingOperateType = typ

	return i
}

func (i *inviteAttendeeService) Id(id string) *inviteAttendeeService {
	i.meetingId = id

	return i
}

func (i *inviteAttendeeService) Attendee(attendee ...*Attendee) *inviteAttendeeService {
	i.attendee = append(i.attendee, attendee...)

	return i
}

func (i *inviteAttendeeService) Body(body *InviteAttendeeReq) *inviteAttendeeService {
	i.body = body

	return i
}

func (i *inviteAttendeeService) BodyString(body string) *inviteAttendeeService {
	i.bodyString = body

	return i
}

func (i *inviteAttendeeService) buildBody() (interface{}, error) {
	if i.body != nil {
		return i.body, nil
	}
	if i.bodyString != "" {
		return i.bodyString, nil
	}
	if len(i.attendee) > 50 {
		i.attendee = i.attendee[:50]
	}

	req := &InviteAttendeeReq{
		MeetingOperateType: i.meetingOperateType,
		MeetingId:          i.meetingId,
		JoinUser:           i.attendee,
	}

	return req, nil
}

func (i *inviteAttendeeService) Do() error {
	req, err := i.buildBody()
	if err != nil {
		return err
	}

	_, err = i.c.performRequest(PerformRequestOptions{
		Method:  http.MethodPost,
		Path:    "/mixapi/inviteSubscribers",
		Params:  nil,
		Body:    req,
		Headers: nil,
	})

	return err
}

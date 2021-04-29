package ysx

import `net/http`

type reInviteService struct {
	c *client

	meetingOperateType meetingOperateType
	meetingId          string
	body               *ReInviteReq
	bodyString         string
}

type ReInviteReq struct {
	MeetingOperateType meetingOperateType `json:"meeting_operate_type"`
	Target             ReInviteTarget     `json:"target"`
}

type ReInviteTarget struct {
	MeetingId string `json:"meeting_id"`
}

func (c *client) ReInvite() *reInviteService {
	return &reInviteService{
		c: c,

		meetingOperateType: MeetingOperateTypePublic,
	}
}

func (r *reInviteService) OperateType(typ meetingOperateType) *reInviteService {
	r.meetingOperateType = typ

	return r
}

func (r *reInviteService) Id(id string) *reInviteService {
	r.meetingId = id

	return r
}

func (r *reInviteService) Body(body *ReInviteReq) *reInviteService {
	r.body = body

	return r
}

func (r *reInviteService) BodyString(body string) *reInviteService {
	r.bodyString = body

	return r
}

func (r *reInviteService) buildBody() (interface{}, error) {
	if r.body != nil {
		return r.body, nil
	}
	if r.bodyString != "" {
		return r.bodyString, nil
	}

	req := &ReInviteReq{
		MeetingOperateType: r.meetingOperateType,
		Target: ReInviteTarget{
			MeetingId: r.meetingId,
		},
	}

	return req, nil
}

func (r *reInviteService) Do() error {
	req, err := r.buildBody()
	if err != nil {
		return err
	}

	_, err = r.c.performRequest(PerformRequestOptions{
		Method:  http.MethodPost,
		Path:    "/mixapi/reInviteSubscribers",
		Params:  nil,
		Body:    req,
		Headers: nil,
	})

	return err
}

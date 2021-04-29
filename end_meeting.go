package ysx

import (
	`net/http`
)

type endMeetingService struct {
	c *client

	meetingOperateType meetingOperateType
	meetingId          string
	body               *EndMeetingReq
	bodyString         string
}

type EndMeetingReq struct {
	MeetingOperateType meetingOperateType `json:"meeting_operate_type"`
	MeetingId          string             `json:"meeting_id"`
}

func (c *client) EndMeeting() *endMeetingService {
	return &endMeetingService{
		c: c,

		meetingOperateType: MeetingOperateTypePublic,
	}
}

func (e *endMeetingService) OperateType(typ meetingOperateType) *endMeetingService {
	e.meetingOperateType = typ

	return e
}

func (e *endMeetingService) Id(id string) *endMeetingService {
	e.meetingId = id

	return e
}

func (e *endMeetingService) Body(body *EndMeetingReq) *endMeetingService {
	e.body = body

	return e
}

func (e *endMeetingService) BodyString(body string) *endMeetingService {
	e.bodyString = body

	return e
}

func (e *endMeetingService) buildBody() (interface{}, error) {
	if e.body != nil {
		return e.body, nil
	}
	if e.bodyString != "" {
		return e.bodyString, nil
	}

	req := &EndMeetingReq{
		MeetingOperateType: e.meetingOperateType,
		MeetingId:          e.meetingId,
	}

	return req, nil
}

func (e *endMeetingService) Do() error {
	req, err := e.buildBody()
	if err != nil {
		return err
	}

	_, err = e.c.performRequest(PerformRequestOptions{
		Method:  http.MethodPost,
		Path:    "/mixapi/endMeeting",
		Params:  nil,
		Body:    req,
		Headers: nil,
	})

	return err
}

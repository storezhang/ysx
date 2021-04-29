package ysx

import (
	`net/http`
	`net/url`
	`strconv`
)

type setChairmanService struct {
	c *client

	meetingOperateType meetingOperateType
	serialNo           string
}

func (c *client) SetChairman() *setChairmanService {
	return &setChairmanService{
		c: c,

		meetingOperateType: MeetingOperateTypePublic,
	}
}

func (r *setChairmanService) OperateType(typ meetingOperateType) *setChairmanService {
	r.meetingOperateType = typ

	return r
}

func (r *setChairmanService) SerialNo(no string) *setChairmanService {
	r.serialNo = no

	return r
}

func (r *setChairmanService) Do() error {
	_, err := r.c.performRequest(PerformRequestOptions{
		Method: http.MethodGet,
		Path:   "/mixapi/setMeetingChairman",
		Params: url.Values{
			"meeting_operate_type": []string{strconv.Itoa(int(r.meetingOperateType))},
			"serialno":             []string{r.serialNo},
		},
		Body:    nil,
		Headers: nil,
	})

	return err
}

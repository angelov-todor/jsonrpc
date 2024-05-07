package service

import "time"

type TimeService struct{}

type TimeRequest struct{}

type TimeResponse struct {
	Time string `json:"time"`
}

func (s *TimeService) GetTime(req *TimeRequest, res *TimeResponse) error {
	res.Time = time.Now().Format(time.RFC3339)
	return nil
}

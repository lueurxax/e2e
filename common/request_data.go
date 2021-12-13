package common

import "time"

// RequestData data object for metrics
type RequestData struct {
	Timestamp time.Time
	Latency   time.Duration
	Code      int
	Method    string
}

// NewRequestData new stage and start timer
func NewRequestData(method string) *RequestData {
	return &RequestData{
		Method:    method,
		Timestamp: time.Now(),
	}
}

// SetProtoCode set code and metering latency
func (s *RequestData) SetProtoCode(code int) {
	s.Latency = time.Since(s.Timestamp)
	s.Code = code
}

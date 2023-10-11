package forum

import "time"

var Sessions = map[string]Session{}

type Session struct {
	Username string
	EndLife  time.Time
}

func (s Session) IsExpired() bool {
	return s.EndLife.Before(time.Now())
}

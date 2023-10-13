package forum

import (
	"time"
)

type Category struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"category"`
	CreationDate time.Time `json:"time"`
}

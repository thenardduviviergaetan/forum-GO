package forum

import (
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	UserType     int64	   `json:"usertype"` // foreign key to UserType table
	UserTypeTxt	 string	   `json:"usertypetxt"`
	Validation   int64	   `json:"validation"` //0 false 1 true
	AskedMod	 int	   `json:"askedmod"`
	CreationDate time.Time `json:"creationtime"`
	FormatedTime string	   `json:"formatedtime"`
}

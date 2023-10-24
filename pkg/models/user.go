package forum

import (
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	UserType     int64     `json:"usertype"` // foreign key to UserType table
	UserTypeTxt  string    `json:"usertypetxt"`
	Validation   int64     `json:"validation"` //0 false 1 true
	AskedMod     int       `json:"askedmod"`
	CreationDate time.Time `json:"creationtime"`
	FormatedTime string    `json:"formatedtime"`
}

type GithubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	// NodeID            string      `json:"node_id"`
	// AvatarURL string `json:"avatar_url"`
	// GravatarID        string      `json:"gravatar_id"`
	// URL               string      `json:"url"`
	// HTMLURL           string `json:"html_url"`
	// FollowersURL      string `json:"followers_url"`
	// FollowingURL      string `json:"following_url"`
	// GistsURL          string `json:"gists_url"`
	// StarredURL        string `json:"starred_url"`
	// SubscriptionsURL  string `json:"subscriptions_url"`
	// OrganizationsURL  string `json:"organizations_url"`
	// ReposURL          string `json:"repos_url"`
	// EventsURL         string `json:"events_url"`
	// ReceivedEventsURL string `json:"received_events_url"`
	// Type              string      `json:"type"`
	// SiteAdmin         bool        `json:"site_admin"`
	// Name              interface{} `json:"name"`
	// Company           interface{} `json:"company"`
	// Blog              string      `json:"blog"`
	// Location          interface{} `json:"location"`
	// Email             interface{} `json:"email"`
	// Hireable          interface{} `json:"hireable"`
	// Bio               interface{} `json:"bio"`
	// TwitterUsername   interface{} `json:"twitter_username"`
	// PublicRepos       int         `json:"public_repos"`
	// PublicGists       int         `json:"public_gists"`
	// Followers         int         `json:"followers"`
	// Following         int         `json:"following"`
	// CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
}

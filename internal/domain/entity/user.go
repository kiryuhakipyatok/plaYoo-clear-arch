package entity

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	Id              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Login           string         `json:"login" gorm:"not null;unique"`
	Telegram      	string         `json:"telegram" gorm:"not null;unique"`
	ChatId          string         `json:"chat_id"`
	Followers       pq.StringArray `gorm:"type:uuid[]" json:"followers"`
	Followings      pq.StringArray `gorm:"type:uuid[]" json:"followings"`
	Rating          float64        `json:"rating"`
	TotalRating     int            `json:"total_rating"`
	NumberOfRatings int            `json:"num_of_ratings"`
	Events          pq.StringArray `gorm:"type:uuid[]" json:"events"`
	Comments        pq.StringArray `gorm:"type:uuid[]" json:"comments"`
	Games           pq.StringArray `gorm:"type:text[]" json:"games"`
	Notifications   pq.StringArray `gorm:"type:uuid[]" json:"notifications"`
	Password        []byte	       `json:"-" gorm:"not null"`
	Discord         string		   `json:"discord"`
	Avatar          string		   `json:"avatar"`
}

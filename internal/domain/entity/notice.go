package entity

import "github.com/google/uuid"

type Notice struct {
	Id      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"notice_id"`
	EventId uuid.UUID `gorm:"type:uuid;not null"`
	Body    string    `json:"body" gorm:"not null"`
}

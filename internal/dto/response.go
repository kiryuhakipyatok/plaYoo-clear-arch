package dto

import (
	"github.com/google/uuid"
	"time"
)

type RegisterResponse struct {
	Id 			uuid.UUID 
	Login 		string
	Telegram 	string
}

type NewsResponse struct {
	Id 			uuid.UUID 
	Title 		string
}

type EventResponse struct {
	Id 			uuid.UUID 
	AuthorId 	uuid.UUID
	Body 		string
	Game		string
	Max			int
	Time		time.Time
}
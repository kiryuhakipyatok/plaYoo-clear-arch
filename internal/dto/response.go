package dto

import "github.com/google/uuid"

type RegisterResponse struct {
	Id 			uuid.UUID 
	Login 		string
	Telegram 	string
}
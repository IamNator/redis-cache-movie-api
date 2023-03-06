package model

import "time"

type (
	Comment struct {
		ID        int       `json:"id"`
		MovieID   int       `json:"movie_id"`
		Comment   string    `json:"comment"`
		IPAddress string    `json:"ip_address"`
		CreatedAt time.Time `json:"created_at"`
	}
)

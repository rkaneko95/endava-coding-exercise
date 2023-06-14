package api

import "time"

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Revoked  Status = "revoked"
)

type Signature struct {
	ID        string    `json:"id"`
	PublicKey string    `json:"publicKey"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}
